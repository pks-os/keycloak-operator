package keycloak

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
)

func TestKeycloakMigration_Test_No_Need_For_Migration_On_Empty_Desired_State(t *testing.T) {
	// given
	migrator := NewDefaultMigrator()
	cr := &v1alpha1.Keycloak{}
	currentState := common.ClusterState{}
	desiredState := common.DesiredClusterState{}

	// when
	migratedActions, error := migrator.Migrate(cr, &currentState, desiredState)

	// then
	assert.Nil(t, error)
	assert.Equal(t, desiredState, migratedActions)
}

func TestKeycloakMigration_Test_No_Need_For_Migration_On_Missing_Deployment_In_Desired_State(t *testing.T) {
	// given
	migrator := NewDefaultMigrator()
	cr := &v1alpha1.Keycloak{}

	keycloakDeployment := model.KeycloakDeployment(cr, nil)
	keycloakDeployment.Spec.Replicas = &[]int32{5}[0]
	keycloakDeployment.Spec.Template.Spec.Containers[0].Image = "old_image" //nolint

	currentState := common.ClusterState{
		KeycloakDeployment: keycloakDeployment,
	}

	desiredState := common.DesiredClusterState{}
	desiredState = append(desiredState, common.GenericUpdateAction{
		Ref: model.KeycloakService(cr),
	})

	// when
	migratedActions, error := migrator.Migrate(cr, &currentState, desiredState)

	// then
	assert.Nil(t, error)
	assert.Equal(t, desiredState, migratedActions)
}

func TestKeycloakMigration_Test_Migrating_Image(t *testing.T) {
	// given
	migrator := NewDefaultMigrator()
	cr := &v1alpha1.Keycloak{}

	keycloakDeployment := model.KeycloakDeployment(cr, model.DatabaseSecret(cr))
	keycloakDeployment.Spec.Replicas = &[]int32{5}[0]
	keycloakDeployment.Spec.Template.Spec.Containers[0].Image = "old_image"

	currentState := common.ClusterState{
		KeycloakDeployment: keycloakDeployment,
	}

	desiredState := common.DesiredClusterState{}
	desiredState = append(desiredState, common.GenericUpdateAction{
		Ref: model.KeycloakDeployment(cr, nil),
	})

	// when
	migratedActions, error := migrator.Migrate(cr, &currentState, desiredState)

	// then
	assert.Nil(t, error)
	assert.Equal(t, int32(0), migratedActions[0].(common.GenericUpdateAction).Ref.(*v1.StatefulSet).Status.Replicas)
}

func TestKeycloakMigration_Test_Migrating_RHSSO_Image(t *testing.T) {
	// given
	migrator := NewDefaultMigrator()
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			Profile: model.RHSSOProfile,
		},
	}

	keycloakDeployment := model.RHSSODeployment(cr, model.DatabaseSecret(cr))
	keycloakDeployment.Spec.Replicas = &[]int32{5}[0]
	keycloakDeployment.Spec.Template.Spec.Containers[0].Image = "old_image"

	currentState := common.ClusterState{
		KeycloakDeployment: keycloakDeployment,
	}

	desiredState := common.DesiredClusterState{}
	desiredState = append(desiredState, common.GenericUpdateAction{
		Ref: model.RHSSODeployment(cr, model.DatabaseSecret(cr)),
	})

	// when
	migratedActions, error := migrator.Migrate(cr, &currentState, desiredState)

	// then
	assert.Nil(t, error)
	assert.Equal(t, int32(0), migratedActions[0].(common.GenericUpdateAction).Ref.(*v1.StatefulSet).Status.Replicas)
}

func TestKeycloakMigration_Test_No_Need_Backup_Without_Migration_Backups_Enabled(t *testing.T) {
	// given
	migrator := NewDefaultMigrator()
	cr := &v1alpha1.Keycloak{}
	cr.Spec.Migration.Backups.Enabled = false

	keycloakDeployment := model.KeycloakDeployment(cr, nil)
	keycloakDeployment.Spec.Template.Spec.Containers[0].Image = "old_image"

	currentState := common.ClusterState{
		KeycloakDeployment: keycloakDeployment,
	}

	desiredState := common.DesiredClusterState{}
	desiredState = append(desiredState, common.GenericUpdateAction{
		Ref: keycloakDeployment,
	})

	// when
	migratedActions, error := migrator.Migrate(cr, &currentState, desiredState)

	// then
	assert.Nil(t, error)
	assert.Equal(t, desiredState, migratedActions)
}

func TestKeycloakMigration_Test_Backup_Happens_With_Migration_Backups_Enabled(t *testing.T) {
	// given
	migrator := NewDefaultMigrator()
	cr := &v1alpha1.Keycloak{}
	cr.Spec.Migration.Backups.Enabled = true

	keycloakDeployment := model.KeycloakDeployment(cr, nil)
	keycloakDeployment.Spec.Template.Spec.Containers[0].Image = "old_image"

	currentState := common.ClusterState{
		KeycloakDeployment: keycloakDeployment,
	}

	desiredState := common.DesiredClusterState{}
	desiredState = append(desiredState, common.GenericUpdateAction{
		Ref: keycloakDeployment,
	})

	// when
	migratedActions, error := migrator.Migrate(cr, &currentState, desiredState)

	// then
	assert.Nil(t, error)
	assert.NotEqual(t, desiredState, migratedActions)
}
