package gamedata

// Constantes de personagem
const (
	// Limites de nome
	MinNameLength = 3
	MaxNameLength = 20

	// Limites de nível e experiência
	StartingLevel = 1
	MaxLevel      = 20
	StartingGold  = 100

	// Limites de atributos
	MinAttributeValue = 8
	MaxAttributeValue = 20
	AttributePoints   = 27

	// Limites de inventário
	MaxInventorySize = 50
	MaxItemQuantity  = 99
	MaxEquipmentSize = 11 // Total de slots de equipamento
)

// Constantes de item
const (
	// Limites de nome do item
	MinItemNameLength = 3
	MaxItemNameLength = 50

	// Limites de descrição do item
	MinDescriptionLength = 10
	MaxDescriptionLength = 200

	// Limites de atributos do item
	MinItemAttributeValue = -5
	MaxItemAttributeValue = 10
)
