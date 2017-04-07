// +build !integration

package module_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/elastic/beats/metricbeat/mb/module"
)

// ExampleWrapper demonstrates how to create a single Wrapper
// from configuration, start the module, and consume events generated by the
// module.
func ExampleWrapper() {
	// Build a configuration object.
	config, err := common.NewConfigFrom(map[string]interface{}{
		"module":     moduleName,
		"metricsets": []string{eventFetcherName},
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create a new Wrapper based on the configuration.
	m, err := module.NewWrapper(config, mb.Registry)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Run the module until done is closed.
	done := make(chan struct{})
	output := m.Start(done)

	// Process events from the output channel until it is closed.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range output {
			// Make rtt a constant so that the output is constant.
			event["metricset"].(common.MapStr)["rtt"] = 111
			fmt.Println(event.StringToPrint())
		}
	}()

	// Simulate running for a while.
	time.Sleep(50 * time.Millisecond)

	// When finished with the module, close the done channel. When the Module
	// stops it will automatically close its output channel so that the output
	// for loop stops.
	close(done)
	wg.Wait()

	// Output:
	// {
	//   "@timestamp": "2016-05-10T23:27:58.485Z",
	//   "_event_metadata": {
	//     "Fields": null,
	//     "FieldsUnderRoot": false,
	//     "Tags": null
	//   },
	//   "fake": {
	//     "eventfetcher": {
	//       "metric": 1
	//     }
	//   },
	//   "metricset": {
	//     "module": "fake",
	//     "name": "eventfetcher",
	//     "rtt": 111
	//   },
	//   "type": "metricsets"
	// }
}

// ExampleRunner demonstrates how to use Runner to start and stop
// a module.
func ExampleRunner() {
	// A *beat.Beat is injected into a Beater when it runs and contains the
	// Publisher used to publish events. This Beat pointer is created here only
	// for demonstration purposes.
	var b *beat.Beat

	config, err := common.NewConfigFrom(map[string]interface{}{
		"module":     moduleName,
		"metricsets": []string{eventFetcherName},
	})
	if err != nil {
		return
	}

	// Create a new Wrapper based on the configuration.
	m, err := module.NewWrapper(config, mb.Registry)
	if err != nil {
		return
	}

	// Create the Runner facade.
	runner := module.NewRunner(b.Publisher.Connect, m)

	// Start the module and have it publish to a new publisher.Client.
	runner.Start()

	// Stop the module. This blocks until all MetricSets in the Module have
	// stopped and the publisher.Client is closed.
	runner.Stop()
}
