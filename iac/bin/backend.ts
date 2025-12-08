import * as cdk from 'aws-cdk-lib';
import {
  BackendStack,
  getProjectName,
  getEnvironment,
  getContainerPort,
  getImageTag,
  getCpu,
  getMemory,
  getDesiredCount,
  getCdkDefaultAccount,
  getCdkDefaultRegion,
  STACK_TYPES,
  createStackName,
  createDefaultTags,
  CloudformationSdkUtils,
  RdsRequests,
  extractRootDomain,
  getBackendDomain
} from 'automation-deploy-template-iac';

(async () => {
    const app = new cdk.App();
    const projectName = app.node.tryGetContext('projectName')
    const environment = getEnvironment(app);
    const rootDomain = app.node.tryGetContext('rootDomain')
    const domainName = rootDomain
    const containerPort = getContainerPort(app);
    const imageTag = getImageTag();
    const cpu = getCpu();
    const memory = getMemory();
    const desiredCount = getDesiredCount();
    const stackName = createStackName(projectName, environment, STACK_TYPES.BACKEND);
    const databaseStackName = createStackName(projectName, environment, STACK_TYPES.DATABASE);

    const databaseStack = await CloudformationSdkUtils.create(databaseStackName);
    const isDatabaseStackDeployed = databaseStack.isDeployed;
    const rdsRequests = isDatabaseStackDeployed ? RdsRequests.build(
        databaseStack.getOutputByKey('ClusterEndpoint'),
        databaseStack.getOutputByKey('ClusterPort'),
        databaseStack.getOutputByKey('DatabaseName'),
        databaseStack.getOutputByKey('SecretArn'),
        databaseStack.getOutputByKey('ClusterArn')
    ) : undefined;

    try {
        // Create Backend Stack using the published npm package
        const backendStack = new BackendStack(app, stackName, {
          projectName,
          environment,
          domainName,
          databaseStackName,
          isDatabaseStackDeployed,
          rdsRequests,
          containerPort,
          imageTag,
          cpu,
          memory,
          desiredCount,
          env: {
            account: getCdkDefaultAccount(),
            region: getCdkDefaultRegion()
          },
          tags: createDefaultTags(projectName, environment, STACK_TYPES.BACKEND, 'standard')
        });

        // Fix for Go backend: Remove hardcoded FastAPI command
        // The automation-deploy-template-iac package hardcodes the command to run 'uvicorn', which fails for Go.
        // We use an escape hatch to remove the command override so the Dockerfile's CMD is used.
        if (backendStack.service) {
          const cfnTaskDef = backendStack.service.taskDefinition.node.defaultChild as cdk.aws_ecs.CfnTaskDefinition;
          // Index 0 corresponds to the backend container (it is the only container in the main task definition)
          // The migration container is in a separate task definition.
          // Explicitly set the command to run the Go binary, overriding the hardcoded 'uvicorn' command.
          cfnTaskDef.addPropertyOverride('ContainerDefinitions.0.Command', ['/main']);
          console.log('✅ Overridden hardcoded FastAPI command with Go binary command');
        }
    
        console.log(`✅ Successfully created ${stackName}`);
      } catch (error) {
        console.error('❌ Failed to create BackendStack:', (error as Error).message);
        process.exit(1);
      }
    })();