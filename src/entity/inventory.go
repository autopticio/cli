package entity

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

type ServiceMetadata struct {
	ServiceName string                 `json:"service_name"`
	Resources   []ResourceMetadata     `json:"resources"`
	MetaData    map[string]interface{} `json:"metadata"`
}

type ResourceMetadata struct {
	ResourceID string                 `json:"resource_id"`
	MetaData   map[string]interface{} `json:"metadata"`
}

// Inventory Entity Commands
func InventoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inventory",
		Short: "Commands related to inventory",
	}

	cmd.AddCommand(makeInventoryCommand())
	// Additional inventory-related commands can be added here

	return cmd
}

func makeInventoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make",
		Short: "Create an inventory file",
		Run: func(cmd *cobra.Command, args []string) {
			out, _ := cmd.Flags().GetString("out")
			log.Printf("Creating inventory at %s\n", out)
			makeInventory(out)
		},
	}
	cmd.Flags().String("out", "", "Output path for the inventory file")

	cmd.MarkFlagRequired("out")
	cmd.MarkFlagFilename("out")
	return cmd
}

func makeInventory(out string) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	// Retrieve AWS account ID
	accountID, err := getAccountID(cfg)
	if err != nil {
		log.Println("Error getting AWS account ID:", err)
		return
	}

	// Retrieve AWS region from config
	region := cfg.Region

	// Get top services by billing
	topServices, err := getTopServicesByCost(cfg)
	if err != nil {
		log.Println("Error getting top services:", err)
		return
	}

	// Iterate over top services and list resources for each
	var results []ServiceMetadata
	for _, service := range topServices {
		serviceCode, ok := mapServiceNameToCode(service)
		if !ok {
			log.Printf("No service code found for service: %s\n", service)
			continue
		}

		// Add metadata: AWS region and account ID
		serviceMeta := ServiceMetadata{
			ServiceName: service,
			MetaData: map[string]interface{}{
				"region":     region,
				"account_id": accountID,
			},
		}

		// List and describe resources for each service based on service code
		switch serviceCode {
		case "ec2":
			resources, err := listEC2Instances(cfg)
			if err != nil {
				log.Println("Error listing EC2 instances:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "s3":
			resources, err := listS3Buckets(cfg)
			if err != nil {
				log.Println("Error listing S3 buckets:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "dynamodb":
			resources, err := listDynamoDBTables(cfg)
			if err != nil {
				log.Println("Error listing DynamoDB tables:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "apigateway":
			// List API Gateway resources
			resources, err := listApiGateways(cfg)
			if err != nil {
				log.Println("Error listing API Gateway resources:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "lambda":
			resources, err := listLambdaFunctions(cfg)
			if err != nil {
				log.Println("Error listing Lambda functions:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "rds":
			resources, err := listRDSInstances(cfg)
			if err != nil {
				log.Println("Error listing RDS instances:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "amazonebs":
			resources, err := listEBSVolumes(cfg)
			if err != nil {
				log.Println("Error listing EBS volumes:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "cloudfront":
			resources, err := listCloudFrontDistributions(cfg)
			if err != nil {
				log.Println("Error listing CloudFront distributions:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "route53":
			resources, err := listRoute53HostedZones(cfg)
			if err != nil {
				log.Println("Error listing Route 53 hosted zones:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "vpc":
			resources, err := listVPCs(cfg)
			if err != nil {
				log.Println("Error listing VPCs:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "elb":
			resources, err := listELBs(cfg)
			if err != nil {
				log.Println("Error listing ELBs:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "cloudwatch":
			resources, err := listCloudWatchMetrics(cfg)
			if err != nil {
				log.Println("Error listing CloudWatch metrics:", err)
				continue
			}
			serviceMeta.Resources = resources

		case "globalaccelerator":
			resources, err := listGlobalAccelerators(cfg)
			if err != nil {
				log.Println("Error listing Global Accelerators:", err)
				continue
			}
			serviceMeta.Resources = resources

		// You can add more cases here for other services

		default:
			log.Printf("Service code not handled: %s\n", serviceCode)
			continue
		}

		results = append(results, serviceMeta)
	}

	// Write results to a JSON file
	writeToJsonFile(results, out)
}

// Get the AWS account ID using STS GetCallerIdentity
func getAccountID(cfg aws.Config) (string, error) {
	svc := sts.NewFromConfig(cfg)
	input := &sts.GetCallerIdentityInput{}
	result, err := svc.GetCallerIdentity(context.TODO(), input)
	if err != nil {
		return "", err
	}

	return *result.Account, nil
}

// Get top services by cost from AWS Cost Explorer
func getTopServicesByCost(cfg aws.Config) ([]string, error) {
	svc := costexplorer.NewFromConfig(cfg)

	endTime := time.Now()
	startTime := endTime.AddDate(0, -1, 0) // Last month

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime.Format("2006-01-02")),
			End:   aws.String(endTime.Format("2006-01-02")),
		},
		Granularity: "MONTHLY",
		Metrics:     []string{"UnblendedCost"},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
		},
	}

	resp, err := svc.GetCostAndUsage(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var topServices []string
	for _, group := range resp.ResultsByTime[0].Groups {
		topServices = append(topServices, group.Keys[0])
	}

	return topServices, nil
}

// Map service names from AWS Cost Explorer or CloudTrail to AWS service codes
func mapServiceNameToCode(serviceName string) (string, bool) {
	serviceMap := map[string]string{
		"Amazon Elastic Compute Cloud - Compute": "ec2",
		"Amazon Simple Storage Service":          "s3",
		"Amazon Relational Database Service":     "rds",
		"Amazon DynamoDB":                        "dynamodb",
		"AWS Lambda":                             "lambda",
		"Amazon CloudFront":                      "cloudfront",
		"Amazon Virtual Private Cloud":           "vpc",
		"Amazon Simple Queue Service":            "sqs",
		"Amazon Simple Notification Service":     "sns",
		"Amazon Elastic Kubernetes Service":      "eks",
		"Amazon Elastic Container Service":       "ecs",
		"Amazon Aurora":                          "aurora",
		"Amazon Redshift":                        "redshift",
		"Amazon Elastic Block Store":             "ebs",
		"AWS Identity and Access Management":     "iam",
		"Amazon Route 53":                        "route53",
		"AmazonCloudWatch":                       "cloudwatch",
		"AWS Key Management Service":             "kms",
		"AWS Glue - Data Integration":            "glue",
		"Amazon SageMaker":                       "sagemaker",
		"AWS Elastic Beanstalk":                  "elasticbeanstalk",
		"AWS Fargate - Serverless Containers":    "fargate",
		"Amazon Elastic File System":             "efs",
		"AWS CloudFormation":                     "cloudformation",
		"AWS Config":                             "config",
		"Amazon Kinesis":                         "kinesis",
		"Amazon API Gateway":                     "apigateway",
		"AWS Step Functions":                     "stepfunctions",
		"Amazon Elastic MapReduce":               "emr",
		"AWS Secrets Manager":                    "secretsmanager",
		"AWS CodeBuild":                          "codebuild",
		"AWS Bedrock":                            "bedrock",
		"Amazon Elastic Load Balancing":          "elb",
		"AWS Global Accelerator":                 "globalaccelerator",
		"Amazon Simple Email Service":            "ses",
		"Amazon Cognito":                         "cognito",
		"EC2 - Other":                            "amazonebs",
	}

	code, exists := serviceMap[serviceName]
	return code, exists
}

// ListGlobalAccelerators retrieves a list of Global Accelerators and their metadata
func listGlobalAccelerators(cfg aws.Config) ([]ResourceMetadata, error) {
	// Set the Global Accelerator endpoint explicitly
	cfg.Region = "us-west-2" // Global Accelerator uses a global endpoint typically aligned with `us-west-2`

	svc := globalaccelerator.NewFromConfig(cfg)
	input := &globalaccelerator.ListAcceleratorsInput{}
	result, err := svc.ListAccelerators(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, accelerator := range result.Accelerators {
		metadata := map[string]interface{}{
			"name":     *accelerator.Name,
			"dns_name": *accelerator.DnsName,
			"status":   accelerator.Status,
			"enabled":  accelerator.Enabled,
			"ip_sets":  accelerator.IpSets,
		}

		resources = append(resources, ResourceMetadata{
			ResourceID: *accelerator.AcceleratorArn,
			MetaData:   metadata,
		})
	}

	return resources, nil
}

// ListCloudWatchMetrics retrieves a list of CloudWatch metrics and their metadata
func listCloudWatchMetrics(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := cloudwatch.NewFromConfig(cfg)
	input := &cloudwatch.ListMetricsInput{}
	result, err := svc.ListMetrics(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, metric := range result.Metrics {
		// Build dimensions information for each metric
		dimensions := make(map[string]string)
		for _, dimension := range metric.Dimensions {
			dimensions[*dimension.Name] = *dimension.Value
		}

		metadata := map[string]interface{}{
			"namespace":   *metric.Namespace,
			"metric_name": *metric.MetricName,
			"dimensions":  dimensions,
		}

		// Using metric name as ResourceID for simplicity
		resources = append(resources, ResourceMetadata{
			ResourceID: *metric.MetricName,
			MetaData:   metadata,
		})
	}

	return resources, nil
}

// ListELBs retrieves a list of Elastic Load Balancers (ELBs) and their metadata
func listELBs(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := elasticloadbalancingv2.NewFromConfig(cfg)
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	result, err := svc.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, elb := range result.LoadBalancers {
		metadata := map[string]interface{}{
			"dns_name":           *elb.DNSName,
			"load_balancer_type": elb.Type,
			"scheme":             elb.Scheme,
			"state":              elb.State.Code,
			"availability_zones": elb.AvailabilityZones,
		}

		resources = append(resources, ResourceMetadata{
			ResourceID: *elb.LoadBalancerArn,
			MetaData:   metadata,
		})
	}

	return resources, nil
}

// ListVPCs retrieves a list of VPCs and their metadata
func listVPCs(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeVpcsInput{}
	result, err := svc.DescribeVpcs(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, vpc := range result.Vpcs {
		metadata := map[string]interface{}{
			"cidr_block": *vpc.CidrBlock,
			"state":      vpc.State,
			"is_default": vpc.IsDefault,
		}

		// Add DNS support and hostnames information, if available
		if vpc.DhcpOptionsId != nil {
			metadata["dhcp_options_id"] = *vpc.DhcpOptionsId
		}

		// Add Tags to metadata, if available
		tags := make(map[string]string)
		for _, tag := range vpc.Tags {
			if tag.Key != nil && tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			}
		}
		metadata["tags"] = tags

		resources = append(resources, ResourceMetadata{
			ResourceID: *vpc.VpcId,
			MetaData:   metadata,
		})
	}

	return resources, nil
}

func listRoute53HostedZones(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := route53.NewFromConfig(cfg)
	input := &route53.ListHostedZonesInput{}
	result, err := svc.ListHostedZones(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, zone := range result.HostedZones {
		metadata := map[string]interface{}{
			"name":                  *zone.Name,
			"resource_record_count": zone.ResourceRecordSetCount,
			"private_zone":          zone.Config.PrivateZone,
		}

		if zone.Config.Comment != nil {
			metadata["comment"] = *zone.Config.Comment
		}

		resources = append(resources, ResourceMetadata{
			ResourceID: *zone.Id,
			MetaData:   metadata,
		})
	}

	return resources, nil
}

// listCloudFrontDistributions retrieves a list of CloudFront distributions and their metadata
func listCloudFrontDistributions(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := cloudfront.NewFromConfig(cfg)
	input := &cloudfront.ListDistributionsInput{}
	result, err := svc.ListDistributions(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	if result.DistributionList != nil {
		for _, distribution := range result.DistributionList.Items {
			metadata := map[string]interface{}{
				"domain_name": distribution.DomainName,
				"status":      distribution.Status,
				"enabled":     distribution.Enabled,
				"origin":      distribution.Origins.Items[0].DomainName, // main origin domain name
			}

			if distribution.Comment != nil {
				metadata["comment"] = *distribution.Comment
			}

			resources = append(resources, ResourceMetadata{
				ResourceID: *distribution.Id,
				MetaData:   metadata,
			})
		}
	}

	return resources, nil
}

// List EC2 instances
func listEC2Instances(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			resources = append(resources, ResourceMetadata{
				ResourceID: *instance.InstanceId,
				MetaData: map[string]interface{}{
					"instance_type": instance.InstanceType,
					"launch_time":   *instance.LaunchTime,
				},
			})
		}
	}

	return resources, nil
}

// ListEBSVolumes retrieves a list of EBS volumes and their metadata
func listEBSVolumes(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeVolumesInput{}
	result, err := svc.DescribeVolumes(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, volume := range result.Volumes {
		metadata := map[string]interface{}{
			"volume_type":       volume.VolumeType,
			"creation_time":     *volume.CreateTime,
			"size_gb":           volume.Size,
			"state":             volume.State,
			"availability_zone": volume.AvailabilityZone,
		}

		if volume.Encrypted != nil {
			metadata["encrypted"] = *volume.Encrypted
		}
		if volume.KmsKeyId != nil {
			metadata["kms_key_id"] = *volume.KmsKeyId
		}

		resources = append(resources, ResourceMetadata{
			ResourceID: *volume.VolumeId,
			MetaData:   metadata,
		})
	}

	return resources, nil
}

// List S3 buckets
func listS3Buckets(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := s3.NewFromConfig(cfg)
	result, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, bucket := range result.Buckets {
		resources = append(resources, ResourceMetadata{
			ResourceID: *bucket.Name,
			MetaData: map[string]interface{}{
				"creation_date": *bucket.CreationDate,
			},
		})
	}

	return resources, nil
}

func writeToJsonFile(data []ServiceMetadata, filename string) {
	// Ensure the directory exists by getting the directory part of the filename
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory:", err)
		return
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write JSON data with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		log.Println("Error writing JSON data:", err)
	}
}

// List DynamoDB tables and describe each
func listDynamoDBTables(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := dynamodb.NewFromConfig(cfg)
	result, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, tableName := range result.TableNames {
		desc, err := svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			log.Printf("Error describing DynamoDB table: %s\n", tableName)
			continue
		}

		resources = append(resources, ResourceMetadata{
			ResourceID: tableName,
			MetaData: map[string]interface{}{
				"table_status": desc.Table.TableStatus,
				"item_count":   *desc.Table.ItemCount,
			},
		})
	}

	return resources, nil
}

// List API Gateway REST APIs
func listApiGateways(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := apigateway.NewFromConfig(cfg)
	result, err := svc.GetRestApis(context.TODO(), &apigateway.GetRestApisInput{})
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, api := range result.Items {
		resources = append(resources, ResourceMetadata{
			ResourceID: *api.Id,
			MetaData: map[string]interface{}{
				"name":        *api.Name,
				"description": *api.Description,
			},
		})
	}

	return resources, nil
}

// List Lambda functions and describe each
func listLambdaFunctions(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := lambda.NewFromConfig(cfg)
	result, err := svc.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{})
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, function := range result.Functions {
		resources = append(resources, ResourceMetadata{
			ResourceID: *function.FunctionName,
			MetaData: map[string]interface{}{
				"runtime":     function.Runtime,
				"last_update": *function.LastModified,
			},
		})
	}

	return resources, nil
}

// List RDS instances and describe each
func listRDSInstances(cfg aws.Config) ([]ResourceMetadata, error) {
	svc := rds.NewFromConfig(cfg)
	result, err := svc.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, err
	}

	var resources []ResourceMetadata
	for _, dbInstance := range result.DBInstances {
		resources = append(resources, ResourceMetadata{
			ResourceID: *dbInstance.DBInstanceIdentifier,
			MetaData: map[string]interface{}{
				"engine":        *dbInstance.Engine,
				"instance_type": *dbInstance.DBInstanceClass,
				"status":        *dbInstance.DBInstanceStatus,
			},
		})
	}

	return resources, nil
}
