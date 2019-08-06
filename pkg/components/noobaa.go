package components

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	extv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	noobaaOperator = "noobaa-operator"
)

// GetNoobaaDeployment returns a Deployment that deploys the nooba-operator
func GetNoobaaDeployment(repository string, tag string, imagePullPolicy string) *appsv1.Deployment {
	registry_name := repository + "/noobaa/" + noobaaOperator
	noobaaImage := fmt.Sprintf("%s:%s", registry_name, tag)
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: noobaaOperator,
			Labels: map[string]string{
				"noobaa-operator": "deployment",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"noobaa-operator": "deployment",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"noobaa-operator": "deployment",
						"app": "noobaa",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "noobaa-operator",
					Containers: []corev1.Container{
						{
							Name:            noobaaOperator,
							Image:           noobaaImage,
							ImagePullPolicy: corev1.PullPolicy(corev1.PullAlways),
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList {
									corev1.ResourceCPU: resource.MustParse("250m"),
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "OPERATOR_NAME",
									Value: "noobaa-operator",
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name: "WATCH_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

// GetNoobaaRole returns the Role required by the nooba-operator to function
func GetNoobaaRole() *rbacv1.Role {
	role := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: noobaaOperator,
			Labels: map[string]string{
				"name": noobaaOperator,
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"noobaa.io",
				},
				Resources: []string{
					"*",
					"noobaas",
					"backingstores",
					"bucketclasses",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
					"services",
					"endpoints",
					"persistentvolumeclaims",
					"events",
					"configmaps",
					"secrets",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"apps",
				},
				Resources: []string{
					"deployments",
					"daemonsets",
					"replicasets",
					"statefulsets",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"monitoring.coreos.com",
				},
				Resources: []string{
					"servicemonitors",
				},
				Verbs: []string{
					"get",
					"create",
				},
			},
			{
				APIGroups: []string{
					"apps",
				},
				ResourceNames: []string{
					"noobaa-operator",
				},
				Resources: []string{
					"deployments/finalizers",
				},
				Verbs: []string{
					"update",
				},
			},
		},
	}
	return role
}

// GetNoobaaClusterRole returns the Role required by the nooba-operatorto function
func GetNoobaaClusterRole() *rbacv1.Role {
	role := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: noobaaOperator,
			Labels: map[string]string{
				"name": noobaaOperator,
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"nodes",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
			{
				APIGroups: []string{
					""},
				Resources: []string{
					"configmaps",
					"secrets",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"namespaces",
				},
				Verbs: []string{
					"get",
				},
			},
			{
				APIGroups: []string{
					"storage.k8s.io",
				},
				Resources: []string{
					"storageclasses",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
			{
				APIGroups: []string{
					"objectbucket.io",
				},
				Resources: []string{
					"*",
				},
				Verbs: []string{
					"*",
				},
			},
		},
	}
	return role
}
// GetNoobaaCRDs returns a list of CustomResourceDefinitions for the Nooba
// custom resources.
func GetNoobaaCRDs() []*extv1beta1.CustomResourceDefinition {
	return []*extv1beta1.CustomResourceDefinition{
		// Noobaa cluster CRD
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "noobaas.noobaa.io",
			},
			Spec: extv1beta1.CustomResourceDefinitionSpec{
				AdditionalPrinterColumns: []extv1beta1.CustomResourceColumnDefinition{
					{
						Name:        "Phase",
						Type:        "string",
						Description: "Phase",
						JSONPath:    ".status.phase",
					},
					{
						Name:        "Mgmt-Endpoints",
						Type:        "string",
						Description: "Mgmt Endpoints",
						JSONPath:    ".status.services.serviceMgmt.nodePorts",
					},
					{
						Name:        "S3-Endpoints",
						Type:        "string",
						Description: "S3 Endpoints",
						JSONPath:    ".status.services.serviceS3.nodePorts",
					},
					{
						Name:        "Image",
						Type:        "string",
						Description: "Actual Image",
						JSONPath:    ".status.actualImage",
					},
					{
						Name:        "Age",
						Type:        "date",
						JSONPath:    ".metadata.creationTimestamp",
					},
				},
				Group:   "noobaa.io",
				Version: "v1alpha1",
				Scope:   "Namespaced",
				Subresources: &extv1beta1.CustomResourceSubresources {
					Status: &extv1beta1.CustomResourceSubresourceStatus{
					},
				},

				Versions: []extv1beta1.CustomResourceDefinitionVersion{
					{
						Name:    "v1alpha1",
						Served:  true,
						Storage: true,
					},
				},

				Names: extv1beta1.CustomResourceDefinitionNames{
					Kind:     "NooBaa",
					ListKind: "NooBaaList",
					Plural:   "noobaas",
					Singular: "noobaa",
					ShortNames: []string{"nb"},
				},

				Validation: &extv1beta1.CustomResourceValidation{
					OpenAPIV3Schema: &extv1beta1.JSONSchemaProps{
						Properties: map[string]extv1beta1.JSONSchemaProps{
							"apiVersion": extv1beta1.JSONSchemaProps{
								Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
								Type: "string",
							},
							"kind": extv1beta1.JSONSchemaProps{
								Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
								Type: "string",
							},
							"metadata": extv1beta1.JSONSchemaProps{
								Description: "Standard object metadata.",
								Type: "object",
							},
							"spec": extv1beta1.JSONSchemaProps{
								Description: "Specification of the desired behavior of the noobaa system.",
								Properties: map[string]extv1beta1.JSONSchemaProps{
									"image": extv1beta1.JSONSchemaProps{
										Description: "Image (optional) overrides the default image for server container",
										Type: "string",
									},
									"imagePullSecret": extv1beta1.JSONSchemaProps{
										Description: "ImagePullSecret (optional) sets a pull secret for the system image",
										Type: "object",
									},
									"mongoImage": extv1beta1.JSONSchemaProps{
										Description: "MongoImage (optional) overrides the default image for mongodb container",
										Type: "string",
									},
									"storageClassName": extv1beta1.JSONSchemaProps{
										Description: "StorageClassName (optional) overrides the default StorageClass for the PVC that the operator creates, this affects where the system stores its database which contains system config, buckets, objects meta-data and mapping file parts to storage locations.",
										Type: "string",
									},
								},
								Type: "object",
							},
							"status": extv1beta1.JSONSchemaProps{
								Description: "Most recently observed status of the noobaa system.",
								Properties: map[string]extv1beta1.JSONSchemaProps{
									"accounts": extv1beta1.JSONSchemaProps{
										Properties: map[string]extv1beta1.JSONSchemaProps{
											"admin": extv1beta1.JSONSchemaProps{
												Properties: map[string]extv1beta1.JSONSchemaProps{
													"secretRef": extv1beta1.JSONSchemaProps{
														Type: "object",
													},
												},
												Required: []string{"secretRef"},
												Type:     "object",
											},
										},
										Required: []string{"admin"},
										Type: "object",
									},
									"actualImage": extv1beta1.JSONSchemaProps{
										Description: "ActualImage is set to report which image the operator is using",
										Type: "string",
									},
									"conditions": extv1beta1.JSONSchemaProps{
										Description: "Current service state of the noobaa system. Based on https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-conditions+patchMergeKey=type+patchStrategy=merge",
										Items: &extv1beta1.JSONSchemaPropsOrArray{
											Schema: &extv1beta1.JSONSchemaProps{
												Properties: map[string]extv1beta1.JSONSchemaProps{
													"lastProbeTime": extv1beta1.JSONSchemaProps{
														Description: "Last time we probed the condition.",
														Format: "date-time",
														Type: "string",
													},
													"lastTransitionTime": extv1beta1.JSONSchemaProps{
														Description: "Last time the condition transitioned from one status to another.",
														Format: "date-time",
														Type:   "string",
													},
													"message": extv1beta1.JSONSchemaProps{
														Description: "Human-readable message indicating details about last transition.",
														Type: "string",
													},
													"reason": extv1beta1.JSONSchemaProps{
														Description: "Unique, one-word, CamelCase reason for the condition's last transition.",
														Type: "string",
													},
													"status": extv1beta1.JSONSchemaProps{
														Description: "Status is the status of the condition.",
														Type: "string",
													},
													"type": extv1beta1.JSONSchemaProps{
														Description: "Type is the type of the condition.",
														Type: "string",
													},
												},
												Required: []string{"type", "status"},
												Type:     "object",
											},
										},
										Type: "array",
									},
									"observedGeneration": extv1beta1.JSONSchemaProps{
										Description: "ObservedGeneration is the most recent generation observed for this noobaa system. It corresponds to the CR generation, which is updated on mutation by the API Server.",
										Format: "int64",
										Type: "integer",
									},
									"phase": extv1beta1.JSONSchemaProps{
										Description: "Phase is a simple, high-level summary of where the System is in its lifecycle",
										Type: "string",
									},
									"readme": extv1beta1.JSONSchemaProps{
										Description: "Readme is a user readable string with explanations on the system",
										Type: "string",
									},
									"services": extv1beta1.JSONSchemaProps{
										Properties: map[string]extv1beta1.JSONSchemaProps{
											"serviceMgmt": extv1beta1.JSONSchemaProps{
												Properties: map[string]extv1beta1.JSONSchemaProps{
													"externalDNS": extv1beta1.JSONSchemaProps{
														Description: "ExternalDNS are external public addresses for the service",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
													"externalIP": extv1beta1.JSONSchemaProps{
														Description: "ExternalIP are external public addresses for the service LoadBalancerPorts such as AWS ELB provide public address and load balancing for the service IngressPorts are manually created public addresses for the service https://kubernetes.io/docs/concepts/services-networking/service/#external-ips https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer https://kubernetes.io/docs/concepts/services-networking/ingress/",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type:   "array",
													},
													"internalDNS": extv1beta1.JSONSchemaProps{
														Description: "InternalDNS are internal addresses of the service inside the cluster",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
													"internalIP": extv1beta1.JSONSchemaProps{
														Description: "InternalIP are internal addresses of the service inside the cluster https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
													"nodePorts": extv1beta1.JSONSchemaProps{
														Description: "NodePorts are the most basic network available it uses the networks available on the hosts of kubernetes nodes. This generally works from within a pod, and from the internal network of the nodes, but may fail from public network. https://kubernetes.io/docs/concepts/services-networking/service/#nodeport",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
													"podPorts": extv1beta1.JSONSchemaProps{
														Description: "PodPorts are the second most basic network address every pod has an IP in the cluster and the pods network is a mesh so the operator running inside a pod in the cluster can use this address. Note: pod IPs are not guaranteed to persist over restarts, so should be rediscovered. Note2: when running the operator outside of the cluster, pod IP is not accessible.",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
												},
												Type: "object",
											},	
											"serviceS3": extv1beta1.JSONSchemaProps{
												Properties: map[string]extv1beta1.JSONSchemaProps{
													"externalDNS": extv1beta1.JSONSchemaProps{
														Description: "ExternalDNS are external public addresses for the service",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
													"externalIP": extv1beta1.JSONSchemaProps{
														Description: "ExternalIP are external public addresses for the service LoadBalancerPorts such as AWS ELB provide public address and load balancing for the service IngressPorts are manually created public addresses for the service https://kubernetes.io/docs/concepts/services-networking/service/#external-ips https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer https://kubernetes.io/docs/concepts/services-networking/ingress/",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type:   "array",
													},
													"internalDNS": extv1beta1.JSONSchemaProps{
														Description: "InternalDNS are internal addresses of the service inside the cluster",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
													"internalIP": extv1beta1.JSONSchemaProps{
														Description: "InternalIP are internal addresses of the service inside the cluster https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
													"nodePorts": extv1beta1.JSONSchemaProps{
														Description: "NodePorts are the most basic network available it uses the networks available on the hosts of kubernetes nodes. This generally works from within a pod, and from the internal network of the nodes, but may fail from public network. https://kubernetes.io/docs/concepts/services-networking/service/#nodeport",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
													"podPorts": extv1beta1.JSONSchemaProps{
														Description: "PodPorts are the second most basic network address every pod has an IP in the cluster and the pods network is a mesh so the operator running inside a pod in the cluster can use this address. Note: pod IPs are not guaranteed to persist over restarts, so should be rediscovered. Note2: when running the operator outside of the cluster, pod IP is not accessible.",
														Items: &extv1beta1.JSONSchemaPropsOrArray{
															Schema: &extv1beta1.JSONSchemaProps{
																Type: "string",
															},
														},
														Type: "array",
													},
												},
												Type: "object",
											},									
										},
										Required: []string{"serviceMgmt", "serviceS3"},
										Type:     "object",
									},
								},
								Required: []string{"observedGeneration", "phase", "actualImage", "accounts", "services", "readme"},
								Type:     "object",
							},
						},
					},
				},
			},
		},
		// BucketClass CRD
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "bucketclasses.noobaa.io",
			},
			Spec: extv1beta1.CustomResourceDefinitionSpec{
				Group:   "noobaa.io",
				Version: "v1alpha1",
				Scope:   "Namespaced",
				Subresources: &extv1beta1.CustomResourceSubresources {
					Status: &extv1beta1.CustomResourceSubresourceStatus{
					},
				},

				Versions: []extv1beta1.CustomResourceDefinitionVersion{
					{
						Name:    "v1alpha1",
						Served:  true,
						Storage: true,
					},
				},

				Names: extv1beta1.CustomResourceDefinitionNames{
					Kind:     "BucketClass",
					ListKind: "BucketClassList",
					Plural:   "bucketclasses",
					Singular: "bucketclass",
				},

				Validation: &extv1beta1.CustomResourceValidation{
					OpenAPIV3Schema: &extv1beta1.JSONSchemaProps{
						Properties: map[string]extv1beta1.JSONSchemaProps{
							"apiVersion": extv1beta1.JSONSchemaProps{
								Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
								Type: "string",
							},
							"kind": extv1beta1.JSONSchemaProps{
								Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
								Type: "string",
							},
							"metadata": extv1beta1.JSONSchemaProps{
								Type: "object",
							},
							"spec": extv1beta1.JSONSchemaProps{
								Type: "object",
							},
							"status": extv1beta1.JSONSchemaProps{
								Type: "object",
							},
						},
					},
				},
			},
		},
		// BackingStores CRD
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "backingstores.noobaa.io",
			},
			Spec: extv1beta1.CustomResourceDefinitionSpec{
				Group:   "noobaa.io",
				Version: "v1alpha1",
				Scope:   "Namespaced",
				Subresources: &extv1beta1.CustomResourceSubresources {
					Status: &extv1beta1.CustomResourceSubresourceStatus{
					},
				},

				Versions: []extv1beta1.CustomResourceDefinitionVersion{
					{
						Name:    "v1alpha1",
						Served:  true,
						Storage: true,
					},
				},

				Names: extv1beta1.CustomResourceDefinitionNames{
					Kind:     "BackingStore",
					ListKind: "BackingStoreList",
					Plural:   "backingstores",
					Singular: "backingstore",
				},

				Validation: &extv1beta1.CustomResourceValidation{
					OpenAPIV3Schema: &extv1beta1.JSONSchemaProps{
						Properties: map[string]extv1beta1.JSONSchemaProps{
							"apiVersion": extv1beta1.JSONSchemaProps{
								Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
								Type: "string",
							},
							"kind": extv1beta1.JSONSchemaProps{
								Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
								Type: "string",
							},
							"metadata": extv1beta1.JSONSchemaProps{
								Type: "object",
							},
							"spec": extv1beta1.JSONSchemaProps{
								Properties: map[string]extv1beta1.JSONSchemaProps{
									"bucketName": extv1beta1.JSONSchemaProps{
										Type: "string",
									},
									"s3Options": extv1beta1.JSONSchemaProps{
										Description: "S3Options specifies client options for the backing store",
										Properties: map[string]extv1beta1.JSONSchemaProps{
											"endpoint": extv1beta1.JSONSchemaProps{
												Description: "Endpoint is the S3 endpoint to use",
												Type: "string",
											},
											"region": extv1beta1.JSONSchemaProps{
												Description: "Region is the AWS region",
												Type:   "string",
											},
											"s3ForcePathStyle": extv1beta1.JSONSchemaProps{
												Description: "S3ForcePathStyle forces the client to send the bucket name in the path aka path-style rather than as a subdomain of the endpoint.",
												Type: "boolean",
											},
											"signatureVersion": extv1beta1.JSONSchemaProps{
												Description: "SignatureVersion specifies the client signature version to use when signing requests.",
												Type: "string",
											},
											"sslDisabled": extv1beta1.JSONSchemaProps{
												Description: "SSLDisabled allows to disable SSL and use plain http",
												Type: "boolean",
											},
										},
										Type: "object",
									},
									"secret": extv1beta1.JSONSchemaProps{
										Description: "Secret refers to a secret that provides the credentials",
										Type: "object",
									},
									"type": extv1beta1.JSONSchemaProps{
										Description: "Type",
										Type: "string",
									},
								},
								Required: []string{"type", "bucketName", "secret"},
								Type:     "object",
							},
							"status": extv1beta1.JSONSchemaProps{
								Properties: map[string]extv1beta1.JSONSchemaProps{
									"conditions": extv1beta1.JSONSchemaProps{
										Description: "Current service state of the noobaa system. Based on https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-conditions+patchMergeKey=type+patchStrategy=merge",
										Items: &extv1beta1.JSONSchemaPropsOrArray{
											Schema: &extv1beta1.JSONSchemaProps{
												Properties: map[string]extv1beta1.JSONSchemaProps{
													"lastProbeTime": extv1beta1.JSONSchemaProps{
														Description: "Last time we probed the condition.",
														Format: "date-time",
														Type: "string",
													},
													"lastTransitionTime": extv1beta1.JSONSchemaProps{
														Description: "Last time the condition transitioned from one status to another.",
														Format: "date-time",
														Type:   "string",
													},
													"message": extv1beta1.JSONSchemaProps{
														Description: "Human-readable message indicating details about last transition.",
														Type: "string",
													},
													"reason": extv1beta1.JSONSchemaProps{
														Description: "Unique, one-word, CamelCase reason for the condition's last transition.",
														Type: "string",
													},
													"status": extv1beta1.JSONSchemaProps{
														Description: "Status is the status of the condition.",
														Type: "string",
													},
													"type": extv1beta1.JSONSchemaProps{
														Description: "Type is the type of the condition.",
														Type: "string",
													},
												},
												Required: []string{"type", "status"},
												Type:     "object",
											},
										},
										Type: "array",
									},
									"phase": extv1beta1.JSONSchemaProps{
										Description: "Phase is a simple, high-level summary of where the System is in its lifecycle",
										Type: "string",
									},
								},
								Required: []string{"phase"},
								Type:     "object",
							},
						},
					},
				},
				AdditionalPrinterColumns: []extv1beta1.CustomResourceColumnDefinition{
					{
						Name:        "Type",
						Type:        "string",
						Description: "Type",
						JSONPath:    ".spec.type",
					},
					{
						Name:        "Bucket-Name",
						Type:        "string",
						Description: "Bucket Name",
						JSONPath:    ".spec.bucketName",
					},
					{
						Name:        "Phase",
						Type:        "string",
						Description: "Phase",
						JSONPath:    ".status.phase",
					},
					{
						Name:        "Age",
						Type:        "date",
						JSONPath:    ".metadata.creationTimestamp",
					},
				},
			},
		},
	}
}

// GetNoobaaCRs is not included as we don't intend for users to directly create NoobaaClusters.
