package pertrie_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
	"github.com/draganm/pertrie"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func init() {
	logger, _ := zap.NewDevelopment()
	if false {

		opts.DefaultContext = logr.NewContext(context.Background(), zapr.NewLogger(logger))
	}
}

var opts = godog.Options{
	Output:        os.Stdout,
	StopOnFailure: true,
	Strict:        true,
	Paths:         []string{"features"},
	NoColors:      true,
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func createTestDatabase(ctx context.Context) (*pertrie.DB, error) {
	td, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("could not create temp dir: %w", err)
	}

	db, err := pertrie.Open(filepath.Join(td, "db"), 0700)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	go func() {
		<-ctx.Done()
		db.Close()
		os.RemoveAll(td)
	}()

	return db, nil

	// pertrie.Open()
}

type State struct {
	db       *pertrie.DB
	lastSize uint64
}

type StateKeyType string

const stateKey = StateKeyType("")

func InitializeScenario(ctx *godog.ScenarioContext) {
	var cancel context.CancelFunc

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ctx, cancel = context.WithCancel(ctx)

		return ctx, nil

	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		cancel()
		return ctx, nil
	})

	state := &State{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {

		// go func() {
		// 	<-ctx.Done()
		// }()

		db, err := createTestDatabase(ctx)
		if err != nil {
			return ctx, err
		}

		state.db = db

		ctx = context.WithValue(ctx, stateKey, state)
		return ctx, nil
	})

	ctx.Step(`^an empty database$`, anEmptyDatabase)
	ctx.Step(`^I get the size of the root$`, iGetTheSizeOfTheRoot)
	ctx.Step(`^the size should be (\d+)$`, theSizeShouldBe)
	ctx.Step(`^I insert a value into the root trie$`, iInsertAValueIntoTheRootTrie)
	ctx.Step(`^the size of the root trie should be (\d+)$`, theSizeOfTheRootTrieShouldBe)
}

func getState(ctx context.Context) *State {
	return ctx.Value(stateKey).(*State)
}

func anEmptyDatabase() error {
	// nothing to be done here, this is a state at beginning of each test
	return nil
}

func iGetTheSizeOfTheRoot(ctx context.Context) error {
	s := getState(ctx)

	return s.db.Write(func(t pertrie.Trie) error {
		s.lastSize = t.Size()
		return nil
	})
}

func theSizeShouldBe(ctx context.Context, expected int64) error {
	s := getState(ctx)
	if s.lastSize != uint64(expected) {
		return fmt.Errorf("expected size to be %d, but was %d", expected, s.lastSize)
	}
	return nil
}

func iInsertAValueIntoTheRootTrie(ctx context.Context) error {
	s := getState(ctx)
	return s.db.Write(func(t pertrie.Trie) error {
		return t.Put([]byte("foo"), []byte("bar"))
	})
}

func theSizeOfTheRootTrieShouldBe(arg1 int) error {
	return godog.ErrPending
}
