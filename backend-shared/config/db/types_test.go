package db

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSetup(t *testing.T) {

	ctx := context.Background()

	// 'testSetup' deletes all database rows that start with 'test-' in the primary key of the row.
	// This ensures a clean slate for the test run.

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()

	applicationStates, err := dbq.UnsafeListAllApplicationStates(ctx)
	assert.NoError(t, err)
	for _, applicationState := range applicationStates {
		if strings.HasPrefix(applicationState.Applicationstate_application_id, "test-") {
			rowsAffected, err := dbq.DeleteApplicationStateById(ctx, applicationState.Applicationstate_application_id)
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, rowsAffected, 1)
			}
		}
	}

	operations, err := dbq.UnsafeListAllOperations(ctx)
	assert.NoError(t, err)
	for _, operation := range operations {

		if strings.HasPrefix(operation.Operation_id, "test-") {
			rowsAffected, err := dbq.DeleteOperationById(ctx, operation.Operation_id, operation.Operation_owner_user_id)
			assert.Equal(t, rowsAffected, 1)
			assert.NoError(t, err)
		}
	}

	applications, err := dbq.UnsafeListAllApplications(ctx)
	assert.NoError(t, err)
	for _, application := range applications {
		if strings.HasPrefix(application.Application_id, "test-") {
			rowsAffected, err := dbq.DeleteApplicationById(ctx, application.Application_id)
			assert.Equal(t, rowsAffected, 1)
			assert.NoError(t, err)
		}
	}

	clusterAccess, err := dbq.UnsafeListAllClusterAccess(ctx)
	assert.NoError(t, err)
	for _, clusterAccess := range clusterAccess {
		if strings.HasPrefix(clusterAccess.Clusteraccess_managed_environment_id, "test-") {
			rowsAffected, err := dbq.DeleteClusterAccessById(ctx, clusterAccess.Clusteraccess_user_id,
				clusterAccess.Clusteraccess_managed_environment_id,
				clusterAccess.Clusteraccess_gitops_engine_instance_id)
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, rowsAffected, 1)
			}
		}
	}

	engineInstances, err := dbq.UnsafeListAllGitopsEngineInstances(ctx)
	assert.NoError(t, err)
	for _, gitopsEngineInstance := range engineInstances {
		if strings.HasPrefix(gitopsEngineInstance.Gitopsengineinstance_id, "test-") {

			rowsAffected, err := dbq.UnsafeDeleteGitopsEngineInstanceById(ctx, gitopsEngineInstance.Gitopsengineinstance_id)

			if !assert.NoError(t, err) {
				return
			}
			if err == nil {
				assert.Equal(t, rowsAffected, 1)
			}
		}
	}

	engineClusters, err := dbq.UnsafeListAllGitopsEngineClusters(ctx)
	assert.NoError(t, err)
	for _, engineCluster := range engineClusters {
		if strings.HasPrefix(engineCluster.Gitopsenginecluster_id, "test-") {
			rowsAffected, err := dbq.AdminDeleteGitopsEngineClusterById(ctx, engineCluster.Gitopsenginecluster_id)
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, rowsAffected, 1)
			}
		}
	}

	managedEnvironments, err := dbq.UnsafeListAllManagedEnvironments(ctx)
	assert.NoError(t, err)
	for _, managedEnvironment := range managedEnvironments {
		if strings.HasPrefix(managedEnvironment.Managedenvironment_id, "test-") {
			rowsAffected, err := dbq.UnsafeDeleteManagedEnvironmentById(ctx, managedEnvironment.Managedenvironment_id)
			assert.Equal(t, rowsAffected, 1)
			assert.NoError(t, err)
		}
	}

	clusterCredentials, err := dbq.UnsafeListAllClusterCredentials(ctx)
	assert.NoError(t, err)
	for _, clusterCredential := range clusterCredentials {
		if strings.HasPrefix(clusterCredential.Clustercredentials_cred_id, "test-") {
			rowsAffected, err := dbq.AdminDeleteClusterCredentialsById(ctx, clusterCredential.Clustercredentials_cred_id)
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, rowsAffected, 1)
			}
		}
	}

	users, err := dbq.UnsafeListAllClusterUsers(ctx)
	assert.NoError(t, err)
	for _, user := range users {
		if strings.HasPrefix(user.Clusteruser_id, "test-") {
			rowsAffected, err := dbq.AdminDeleteClusterUserById(ctx, (user.Clusteruser_id))
			assert.Equal(t, rowsAffected, 1)
			assert.NoError(t, err)
		}
	}

	err = dbq.CreateClusterUser(ctx, testClusterUser)
	assert.NoError(t, err)

}

func testTeardown(t *testing.T) {
	// Currently unused
}

// Ensure that the we are able to select on all the fields of the database.
func TestSelectOnAllTables(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)
	ctx := context.Background()

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()

	_, err = dbq.UnsafeListAllApplicationStates(ctx)
	assert.NoError(t, err)

	_, err = dbq.UnsafeListAllApplications(ctx)
	assert.NoError(t, err)

	_, err = dbq.UnsafeListAllClusterAccess(ctx)
	assert.NoError(t, err)

	_, err = dbq.UnsafeListAllClusterCredentials(ctx)
	assert.NoError(t, err)

	_, err = dbq.UnsafeListAllClusterUsers(ctx)
	assert.NoError(t, err)

	_, err = dbq.UnsafeListAllGitopsEngineClusters(ctx)
	assert.NoError(t, err)

	_, err = dbq.UnsafeListAllGitopsEngineInstances(ctx)
	assert.NoError(t, err)

	_, err = dbq.UnsafeListAllManagedEnvironments(ctx)
	assert.NoError(t, err)

	_, err = dbq.UnsafeListAllOperations(ctx)
	assert.NoError(t, err)

}

func TestCreateApplication(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()

	ctx := context.Background()

	managedEnv, _, engineInstance, clusterAccess, err := createSampleData(t, dbq)
	if !assert.NoError(t, err) {
		return
	}

	application := &Application{
		Application_id:          "test-my-application",
		Name:                    "my-application",
		Spec_field:              "{}",
		Engine_instance_inst_id: engineInstance.Gitopsengineinstance_id,
		Managed_environment_id:  managedEnv.Managedenvironment_id,
	}

	err = dbq.CreateApplication(ctx, application, clusterAccess.Clusteraccess_user_id)
	if !assert.NoError(t, err) {
		return
	}

	applicationRes, err := dbq.UnsafeGetApplicationById(ctx, application.Application_id)
	if !assert.Equal(t, application.Application_id, applicationRes.Application_id) {
		return
	}
	if !assert.NoError(t, err) {
		return
	}

	rowsAffected, err := dbq.DeleteApplicationById(ctx, application.Application_id)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, rowsAffected, 1) {
		return
	}

	applicationRes, err = dbq.UnsafeGetApplicationById(ctx, application.Application_id)
	if !assert.Nil(t, applicationRes) {
		return
	}
	if !assert.Error(t, err) {
		return
	}

}

func TestDeploymentToApplicationMapping(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)
	ctx := context.Background()

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()

	res, err := dbq.GetDeploymentToApplicationMappingById(ctx, "meow")
	fmt.Println(res, err)

}

func TestGitopsEngineInstanceAndCluster(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()
	ctx := context.Background()

	_, engineCluster, engineInstance, _, err := createSampleData(t, dbq)
	if !assert.NoError(t, err) {
		return
	}

	result, err := dbq.GetGitopsEngineClusterById(ctx, engineCluster.Gitopsenginecluster_id, testClusterUser.Clusteruser_id)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, engineCluster, result) {
		return
	}

	rowsAffected, err := dbq.DeleteGitopsEngineInstanceById(ctx, engineInstance.Gitopsengineinstance_id, testClusterUser.Clusteruser_id)
	assert.Equal(t, rowsAffected, 1)
	assert.NoError(t, err)

	// Get should return no found, after the delete
	_, err = dbq.GetGitopsEngineInstanceById(ctx, engineCluster.Gitopsenginecluster_id, testClusterUser.Clusteruser_id)
	if !assert.Error(t, err) {
		return
	}
	assert.True(t, IsResultNotFoundError(err))

	rowsAffected, err = dbq.AdminDeleteGitopsEngineClusterById(ctx, engineCluster.Gitopsenginecluster_id)
	assert.Equal(t, rowsAffected, 1)
	assert.NoError(t, err)

	_, err = dbq.GetGitopsEngineClusterById(ctx, engineCluster.Gitopsenginecluster_id, testClusterUser.Clusteruser_id)
	assert.Error(t, err)
	assert.True(t, IsResultNotFoundError(err))
}

func createSampleData(t *testing.T, dbq AllDatabaseQueries) (*ManagedEnvironment, *GitopsEngineCluster, *GitopsEngineInstance, *ClusterAccess, error) {

	ctx := context.Background()

	var err error

	managedEnvironment, engineCluster, engineInstance, clusterAccess := generateSampleData()

	if err = dbq.CreateManagedEnvironment(ctx, &managedEnvironment); err != nil {
		return nil, nil, nil, nil, err
	}

	if err = dbq.CreateClusterAccess(ctx, &clusterAccess); err != nil {
		return nil, nil, nil, nil, err
	}

	if err = dbq.CreateGitopsEngineCluster(ctx, &engineCluster); err != nil {
		return nil, nil, nil, nil, err
	}

	if err = dbq.CreateGitopsEngineInstance(ctx, &engineInstance); err != nil {
		return nil, nil, nil, nil, err
	}

	return &managedEnvironment, &engineCluster, &engineInstance, &clusterAccess, nil

}

func TestManagedEnvironment(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)
	ctx := context.Background()

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()

	managedEnvironment, _, _, _, err := createSampleData(t, dbq)
	if !assert.NoError(t, err) {
		return
	}

	result, err := dbq.GetManagedEnvironmentById(ctx, managedEnvironment.Managedenvironment_id, testClusterUser.Clusteruser_id)
	assert.NoError(t, err)
	if err == nil {
		assert.Equal(t, managedEnvironment.Name, result.Name)
	}

	_, err = dbq.GetManagedEnvironmentById(ctx, managedEnvironment.Managedenvironment_id, "another-user")
	assert.NotNil(t, err) // deleting from another user should fail
	assert.True(t, IsResultNotFoundError(err))

	rowsAffected, _ := dbq.DeleteManagedEnvironmentById(ctx, managedEnvironment.Managedenvironment_id, "another-user")
	assert.Equal(t, rowsAffected, 0) // deleting from another user should not return any rows

	rowsAffected, err = dbq.DeleteManagedEnvironmentById(ctx, managedEnvironment.Managedenvironment_id, testClusterUser.Clusteruser_id)
	assert.Equal(t, rowsAffected, 1)
	assert.NoError(t, err)

	_, err = dbq.GetManagedEnvironmentById(ctx, managedEnvironment.Managedenvironment_id, testClusterUser.Clusteruser_id)
	assert.NotNil(t, err)
	assert.True(t, IsResultNotFoundError(err))

}

func TestOperation(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()

	_, _, engineInstance, _, err := createSampleData(t, dbq)
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()

	operation := &Operation{
		Instance_id:             engineInstance.Gitopsengineinstance_id,
		Resource_id:             "fake resource id",
		Resource_type:           "GitopsEngineInstance",
		State:                   OperationState_Waiting,
		Operation_owner_user_id: testClusterUser.Clusteruser_id,
	}

	err = dbq.CreateOperation(ctx, operation, "some-other-user")
	assert.Error(t, err)

	err = dbq.CreateOperation(ctx, operation, operation.Operation_owner_user_id)
	assert.NoError(t, err)

	result, err := dbq.GetOperationById(ctx, operation.Operation_id, operation.Operation_owner_user_id)
	assert.NoError(t, err)
	assert.Equal(t, result.Operation_id, operation.Operation_id)

	_, err = dbq.GetOperationById(ctx, operation.Operation_id, "another-user")
	if !assert.Error(t, err) {
		return
	}
	assert.True(t, IsResultNotFoundError(err))

	rowsAffected, _ := dbq.DeleteOperationById(ctx, operation.Operation_id, "another-user")
	assert.Equal(t, rowsAffected, 0)

	rowsAffected, err = dbq.DeleteOperationById(ctx, operation.Operation_id, operation.Operation_owner_user_id)
	assert.Equal(t, rowsAffected, 1)
	assert.NoError(t, err)
}

func TestClusterUser(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)
	ctx := context.Background()

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()

	clusterUser := ClusterUser{
		Clusteruser_id: "test-my-cluster-user-2",
		User_name:      "cluster-mccluster",
	}
	err = dbq.CreateClusterUser(ctx, &clusterUser)
	assert.NoError(t, err)

	retrievedClusterUser, err := dbq.GetClusterUserById(ctx, clusterUser.Clusteruser_id)
	assert.NoError(t, err)
	assert.Equal(t, clusterUser.User_name, retrievedClusterUser.User_name)

	rowsAffected, err := dbq.AdminDeleteClusterUserById(ctx, clusterUser.Clusteruser_id)
	assert.Equal(t, rowsAffected, 1)
	assert.NoError(t, err)

	_, err = dbq.GetClusterUserById(ctx, clusterUser.Clusteruser_id)
	if !assert.Error(t, err) {
		return
	}
	assert.True(t, IsResultNotFoundError(err))

	_, err = dbq.GetClusterUserById(ctx, "does-not-exist")
	if !assert.Error(t, err) {
		return
	}
	assert.True(t, IsResultNotFoundError(err))

}

func TestClusterCredentials(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	dbq, err := NewUnsafePostgresDBQueries(true, true)
	if !assert.NoError(t, err) {
		return
	}
	defer dbq.CloseDatabase()
	ctx := context.Background()

	clusterCredentials := ClusterCredentials{
		Clustercredentials_cred_id:  "test-cluster-creds-test",
		Host:                        "host",
		Kube_config:                 "kube-config",
		Kube_config_context:         "kube-config-context",
		Serviceaccount_bearer_token: "serviceaccount_bearer_token",
		Serviceaccount_ns:           "Serviceaccount_ns",
	}

	err = dbq.CreateClusterCredentials(ctx, &clusterCredentials)
	assert.NoError(t, err)

	// Create managed environment, and cluster access, so the non-unsafe get works below
	{

		managedEnvironment := ManagedEnvironment{
			Managedenvironment_id: "test-managed-env-914",
			Clustercredentials_id: clusterCredentials.Clustercredentials_cred_id,
			Name:                  "my env",
		}
		err = dbq.CreateManagedEnvironment(ctx, &managedEnvironment)
		if !assert.NoError(t, err) {
			return
		}

		clusterAccess := ClusterAccess{
			Clusteraccess_user_id:                   testClusterUser.Clusteruser_id,
			Clusteraccess_managed_environment_id:    managedEnvironment.Managedenvironment_id,
			Clusteraccess_gitops_engine_instance_id: "fake-engine-instance-id",
		}

		err = dbq.CreateClusterAccess(ctx, &clusterAccess)
		if !assert.NoError(t, err) {
			return
		}
	}

	retrievedClusterCredentials, err := dbq.UnsafeGetClusterCredentialsById(ctx, clusterCredentials.Clustercredentials_cred_id)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, clusterCredentials.Host, retrievedClusterCredentials.Host)
	assert.Equal(t, clusterCredentials.Kube_config, retrievedClusterCredentials.Kube_config)
	assert.Equal(t, clusterCredentials.Kube_config_context, retrievedClusterCredentials.Kube_config_context)

	retrievedClusterCredentials, err = dbq.GetClusterCredentialsById(ctx, clusterCredentials.Clustercredentials_cred_id, testClusterUser.Clusteruser_id)
	if !assert.NoError(t, err) ||
		!assert.NotNil(t, retrievedClusterCredentials) {
		return
	}

	assert.Equal(t, clusterCredentials.Host, retrievedClusterCredentials.Host)
	assert.Equal(t, clusterCredentials.Kube_config, retrievedClusterCredentials.Kube_config)
	assert.Equal(t, clusterCredentials.Kube_config_context, retrievedClusterCredentials.Kube_config_context)

	rowsAffected, err := dbq.AdminDeleteClusterCredentialsById(ctx, clusterCredentials.Clustercredentials_cred_id)
	assert.Equal(t, rowsAffected, 1)
	assert.NoError(t, err)

	_, err = dbq.UnsafeGetClusterCredentialsById(ctx, clusterCredentials.Clustercredentials_cred_id)
	if !assert.Error(t, err) {
		return
	}
	assert.True(t, IsResultNotFoundError(err))

}

var testClusterUser = &ClusterUser{
	Clusteruser_id: "test-user",
	User_name:      "test-user",
}

func generateSampleData() (ManagedEnvironment, GitopsEngineCluster, GitopsEngineInstance, ClusterAccess) {

	managedEnvironment := ManagedEnvironment{
		Managedenvironment_id: "test-managed-environment",
		Name:                  "my managed environment",
		Clustercredentials_id: "test-cluster-creds",
	}

	engineCluster := GitopsEngineCluster{
		Gitopsenginecluster_id: "test-engine-cluster-id",
		Clustercredentials_id:  "test-fake-creds-id",
	}

	engineInstance := GitopsEngineInstance{
		Gitopsengineinstance_id: "test-fake-engine-instance",
		Namespace_name:          "my-namespace",
		Namespace_uid:           "my-namespace-uid",
		EngineCluster_id:        engineCluster.Gitopsenginecluster_id,
	}

	clusterAccess := ClusterAccess{
		Clusteraccess_user_id:                   testClusterUser.Clusteruser_id,
		Clusteraccess_managed_environment_id:    managedEnvironment.Managedenvironment_id,
		Clusteraccess_gitops_engine_instance_id: engineInstance.Gitopsengineinstance_id,
	}

	return managedEnvironment, engineCluster, engineInstance, clusterAccess
}
