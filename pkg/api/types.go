package api

// TypeMeta describes the type of object
type TypeMeta struct {
	Kind string
}

// ObjectMeta describes the metadata of object
type ObjectMeta struct {
	Name   string
	Labels map[string]string
}

// Job is a batch job resource
type Job struct {
	TypeMeta
	ObjectMeta
	Spec JobSpec
}

// JobSpec is a describtion of Job
type JobSpec struct {
	Containers    []Container
	NodeName      string
	NodeSelectors map[string]string
}

// Container is a describtion of docker container
type Container struct {
	Image     string
	Name      string
	Command   []string
	Resources ResourceRequirement
}

// ResourceRequirement describes limit and required resources of a contianer
type ResourceRequirement struct {
	Limit    ResourceList
	Required ResourceList
}

// ResourceList describes a list of resource, i.g. CPU, Memory.
type ResourceList struct {
	CPU    int
	Memory int
}

// Node is a os node resource
type Node struct {
	TypeMeta
	ObjectMeta
	Spec NodeSpec
}

// NodeSpec is a describtion of node
type NodeSpec struct {
	ExternalIP string
}
