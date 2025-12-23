package listener

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/events"
)

type TaskEventListener struct{}

func NewTaskEventListener() *TaskEventListener {
	return &TaskEventListener{}
}

func (l *TaskEventListener) HandleTaskAssigned(payload []byte) error {
	var event events.TaskAssignedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskAssignedEvent: %w", err)
	}

	log.Printf("🎯 YENİ TASK ATAMASI!")
	log.Printf("   👤 Kullanıcı: %s (%s)", event.UserName, event.UserEmail)
	log.Printf("   📋 Task: %s (ID: %s)", event.TaskTitle, event.TaskID)
	log.Printf("   📧 Email gönderiliyor...")

	if err := l.sendEmail(event); err != nil {
		log.Printf("   ❌ Email gönderilemedi: %v", err)
		return err
	}

	log.Printf("   ✅ Email başarıyla gönderildi!")
	return nil
}

func (l *TaskEventListener) sendEmail(event events.TaskAssignedEvent) error {
	log.Printf("   📨 TO: %s", event.UserEmail)
	log.Printf("   📨 SUBJECT: Yeni Görev Atandı: %s", event.TaskTitle)
	log.Printf("   📨 BODY: Merhaba %s, size yeni bir görev atandı!", event.UserName)

	return nil
}
