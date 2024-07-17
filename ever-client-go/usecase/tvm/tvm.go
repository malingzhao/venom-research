package tvm

import (
	"github.com/markgenuine/ever-client-go/domain"
)

type tvm struct {
	config domain.ClientConfig
	client domain.ClientGateway
}

func NewTvm(
	config domain.ClientConfig,
	client domain.ClientGateway,
) domain.TvmUseCase {
	return &tvm{
		config: config,
		client: client,
	}
}

// RunExecutor - Emulates all the phases of contract execution locally.
// Performs all the phases of contract execution on Transaction Executor
// - the same component that is used on Validator Nodes.
func (t *tvm) RunExecutor(pORE *domain.ParamsOfRunExecutor) (*domain.ResultOfRunExecuteMessage, error) {
	result := new(domain.ResultOfRunExecuteMessage)
	err := t.client.GetResult("tvm.run_executor", pORE, result)
	return result, err
}

// RunTvm - Executes get-methods of ABI-compatible contracts.
func (t *tvm) RunTvm(pORT *domain.ParamsOfRunTvm) (*domain.ResultOfRunTvm, error) {
	result := new(domain.ResultOfRunTvm)
	err := t.client.GetResult("tvm.run_tvm", pORT, result)
	return result, err
}

// RunGet - Executes a get-method of FIFT contract.
func (t *tvm) RunGet(pORG *domain.ParamsOfRunGet) (*domain.ResultOfRunGet, error) {
	result := new(domain.ResultOfRunGet)
	err := t.client.GetResult("tvm.run_get", pORG, result)
	return result, err
}
