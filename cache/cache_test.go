package cache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	t.Run("создание_cache", func(t *testing.T) {
		c := NewCache[string](time.Second)
		if c == nil {
			t.Fatal("ожидали cache")
		}
		if c.items == nil {
			t.Fatal("ожидали инициализированную map")
		}
	})
}

func TestCache_SetAndGet(t *testing.T) {
	t.Run("ключ_существует", func(t *testing.T) {
		c := NewCache[string](time.Second)
		c.Set("k", "v")

		got, ok := c.Get("k")
		if !ok {
			t.Fatal("ожидали значение")
		}
		if got != "v" {
			t.Fatalf("получили %q, ожидали v", got)
		}
	})

	t.Run("ключ_отсутствует", func(t *testing.T) {
		c := NewCache[string](time.Second)

		_, ok := c.Get("missing")
		if ok {
			t.Fatal("ожидали отсутствие значения")
		}
	})
}

func TestCache_GetExpired(t *testing.T) {
	t.Run("ttl_истёк", func(t *testing.T) {
		c := NewCache[string](5 * time.Millisecond)
		c.Set("k", "v")
		time.Sleep(20 * time.Millisecond)

		_, ok := c.Get("k")
		if ok {
			t.Fatal("ожидали отсутствие значения")
		}
	})
}

func TestCache_Evict(t *testing.T) {
	t.Run("удаление_по_ключу", func(t *testing.T) {
		c := NewCache[string](time.Second)
		c.Set("k", "v")
		c.Evict("k")

		_, ok := c.Get("k")
		if ok {
			t.Fatal("ожидали отсутствие значения")
		}
	})
}

func TestCache_Cleanup(t *testing.T) {
	t.Run("очистка_истёкших_записей", func(t *testing.T) {
		c := NewCache[string](5 * time.Millisecond)
		c.Set("expired", "v1")
		time.Sleep(20 * time.Millisecond)
		c.Set("alive", "v2")

		c.Cleanup()

		if _, ok := c.items["expired"]; ok {
			t.Fatal("ожидали удаление истёкшей записи")
		}
		if _, ok := c.items["alive"]; !ok {
			t.Fatal("ожидали сохранение актуальной записи")
		}
	})
}
