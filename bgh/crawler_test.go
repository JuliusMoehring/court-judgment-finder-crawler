package bgh

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_links(t *testing.T) {
	t.Run("Correctly adds links concurrently", func(t *testing.T) {
		links := newLinks()

		var wg sync.WaitGroup
		numGoroutines := 1000

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				links.addLink(fmt.Sprintf("https://example.com/%d", index))
			}(i)
		}

		wg.Wait()

		assert.Equal(t, numGoroutines, len(links.getLinks()), "Should have added 1000 links")
	})

	t.Run("Correctly adds duplicate links", func(t *testing.T) {
		links := newLinks()

		var wg sync.WaitGroup
		numGoroutines := 1000

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				links.addLink("https://example.com")
			}()
		}

		wg.Wait()

		assert.Equal(t, 1, len(links.getLinks()), "Should have added 1 link")
	})
}
