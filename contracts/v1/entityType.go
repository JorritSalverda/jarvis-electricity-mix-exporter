package contracts

type EntityType string

const (
	EntityType_ENTITY_TYPE_INVALID EntityType = ""
	EntityType_ENTITY_TYPE_TARIFF  EntityType = "ENTITY_TYPE_TARIFF"
	EntityType_ENTITY_TYPE_ZONE    EntityType = "ENTITY_TYPE_ZONE"
	EntityType_ENTITY_TYPE_DEVICE  EntityType = "ENTITY_TYPE_DEVICE"
)
