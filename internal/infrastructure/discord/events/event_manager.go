package events

import (
	"fmt"
	"log"
	"sync"

	"sirdraith/internal/domain/repository"
	"sirdraith/internal/infrastructure/discord/commands"

	"github.com/bwmarrin/discordgo"
)

// EventType representa o tipo de evento do Discord
type EventType string

const (
	EventReady             EventType = "ready"
	EventGuildCreate       EventType = "guild_create"
	EventGuildDelete       EventType = "guild_delete"
	EventGuildMemberAdd    EventType = "guild_member_add"
	EventGuildMemberRem    EventType = "guild_member_remove"
	EventMessageCreate     EventType = "message_create"
	EventMessageDelete     EventType = "message_delete"
	EventMessageUpdate     EventType = "message_update"
	EventInteractionCreate EventType = "interaction_create"
)

// EventHandler é a interface que todos os handlers de eventos devem implementar
type EventHandler interface {
	Handle(s *discordgo.Session, i interface{}) error
}

// EventHandlerFunc é uma função que implementa EventHandler
type EventHandlerFunc func(s *discordgo.Session, i interface{}) error

// Handle implementa EventHandler
func (f EventHandlerFunc) Handle(s *discordgo.Session, i interface{}) error {
	return f(s, i)
}

// EventManager gerencia os handlers de eventos do Discord
type EventManager struct {
	handlers        map[EventType][]EventHandler
	configRepo      repository.ConfigRepository
	commandRegistry *commands.CommandRegistry
	session         *discordgo.Session
	mu              sync.RWMutex
}

// NewEventManager cria uma nova instância do EventManager
func NewEventManager(session *discordgo.Session, configRepo repository.ConfigRepository, registry *commands.CommandRegistry) *EventManager {
	return &EventManager{
		handlers:        make(map[EventType][]EventHandler),
		configRepo:      configRepo,
		commandRegistry: registry,
		session:         session,
	}
}

// AddHandler adiciona um novo handler para um tipo de evento
func (em *EventManager) AddHandler(eventType EventType, handler EventHandler) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if em.handlers[eventType] == nil {
		em.handlers[eventType] = make([]EventHandler, 0)
	}
	em.handlers[eventType] = append(em.handlers[eventType], handler)
}

// HandleEvent processa um evento usando os handlers registrados
func (em *EventManager) HandleEvent(eventType EventType, s *discordgo.Session, i interface{}) error {
	em.mu.RLock()
	handlers := em.handlers[eventType]
	em.mu.RUnlock()

	if len(handlers) == 0 {
		return fmt.Errorf("nenhum handler registrado para o evento %s", eventType)
	}

	var lastErr error
	for _, handler := range handlers {
		if err := handler.Handle(s, i); err != nil {
			log.Printf("Erro ao processar evento %s: %v", eventType, err)
			lastErr = err
		}
	}

	return lastErr
}

// WrapHandler cria um wrapper para o handler do discordgo que direciona o evento para o EventManager
func WrapHandler(eventType EventType, em *EventManager) interface{} {
	switch eventType {
	case EventReady:
		return func(s *discordgo.Session, i *discordgo.Ready) {
			_ = em.HandleEvent(eventType, s, i)
		}
	case EventGuildCreate:
		return func(s *discordgo.Session, i *discordgo.GuildCreate) {
			_ = em.HandleEvent(eventType, s, i)
		}
	case EventGuildDelete:
		return func(s *discordgo.Session, i *discordgo.GuildDelete) {
			_ = em.HandleEvent(eventType, s, i)
		}
	case EventGuildMemberAdd:
		return func(s *discordgo.Session, i *discordgo.GuildMemberAdd) {
			_ = em.HandleEvent(eventType, s, i)
		}
	case EventGuildMemberRem:
		return func(s *discordgo.Session, i *discordgo.GuildMemberRemove) {
			_ = em.HandleEvent(eventType, s, i)
		}
	case EventMessageCreate:
		return func(s *discordgo.Session, i *discordgo.MessageCreate) {
			_ = em.HandleEvent(eventType, s, i)
		}
	case EventMessageDelete:
		return func(s *discordgo.Session, i *discordgo.MessageDelete) {
			_ = em.HandleEvent(eventType, s, i)
		}
	case EventMessageUpdate:
		return func(s *discordgo.Session, i *discordgo.MessageUpdate) {
			_ = em.HandleEvent(eventType, s, i)
		}
	case EventInteractionCreate:
		return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_ = em.HandleEvent(eventType, s, i)
		}
	default:
		log.Printf("Tipo de evento não suportado: %s", eventType)
		return nil
	}
}

// RegisterDefaultHandlers registra os handlers padrão para os eventos
func (em *EventManager) RegisterDefaultHandlers() {
	// Ready Handler
	em.AddHandler(EventReady, NewReadyHandler(em.session))

	// Guild Handlers
	em.AddHandler(EventGuildCreate, NewGuildCreateHandler(em.session))
	em.AddHandler(EventGuildDelete, NewGuildDeleteHandler())
	em.AddHandler(EventGuildMemberAdd, NewGuildMemberAddHandler(em.configRepo))
	em.AddHandler(EventGuildMemberRem, NewGuildMemberRemoveHandler(em.configRepo))

	// Message Handlers
	messageHandler := NewMessageHandler()
	em.AddHandler(EventMessageCreate, messageHandler)
	em.AddHandler(EventMessageDelete, messageHandler)
	em.AddHandler(EventMessageUpdate, messageHandler)

	// Interaction Handler
	em.AddHandler(EventInteractionCreate, NewInteractionCreateHandler())

	log.Println("Handlers de eventos registrados com sucesso!")
}

// HandleInteractionCreate processa eventos de interação
func (em *EventManager) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Processa a interação através do registro de comandos
	if err := em.commandRegistry.HandleInteraction(i); err != nil {
		log.Printf("Erro ao processar interação: %v", err)
	}
}
