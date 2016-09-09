package layerx_brain_client

type BrainAssignTasksMessage struct {
	NodeId string `json:"node_id"`
	TaskIds []string `json:"task_ids"`
}

type MigrateTaskMessage struct {
	DestinationNodeId string `json:"destination_node_id"`
	TaskIds []string `json:"task_ids"`
}
