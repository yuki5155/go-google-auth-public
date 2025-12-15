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
  ContainerConfigRequests,
  extractRootDomain,
  getBackendDomain,
  getFrontendDomain
} from 'automation-deploy-template-iac';

(async () => {
    const app = new cdk.App();
    const projectName = app.node.tryGetContext('projectName')
    const environment = getEnvironment(app);
    const rootDomain = app.node.tryGetContext('rootDomain')
    const domainName = rootDomain ? getBackendDomain(extractRootDomain(rootDomain), environment) : undefined;
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

    // Setup environment variables for CORS
    const frontendDomain = rootDomain ? getFrontendDomain(extractRootDomain(rootDomain), environment) : undefined;
    const frontendUrl = frontendDomain ? `https://${frontendDomain}` : 'http://localhost:5173';
    
    console.log('=== Backend Environment Configuration ===');
    console.log(`Frontend URL: ${frontendUrl}`);
    console.log(`Backend Domain: ${domainName || 'Not specified'}`);
    console.log(`Environment: ${environment}`);
    
    const containerConfigRequests = ContainerConfigRequests.build({
      ALLOWED_ORIGINS: frontendUrl,
      FRONTEND_URL: frontendUrl,
      GO_ENV: environment === 'prod' ? 'production' : 'development',
      PORT: containerPort.toString()
    });

    try {
        // Create Backend Stack using the published npm package
        const backendStack = new BackendStack(app, stackName, {
          projectName,
          environment,
          domainName,
          databaseStackName,
          isDatabaseStackDeployed,
          rdsRequests,
          containerConfigRequests,
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
    
        console.log(`✅ Successfully created ${stackName}`);
      } catch (error) {
        console.error('❌ Failed to create BackendStack:', (error as Error).message);
        process.exit(1);
      }
    })();