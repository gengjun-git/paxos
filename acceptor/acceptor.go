package acceptor

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"
)

const epochFileName = "epoch"
const proposalFileName = "proposal"
const fileReadBufferSize = 4096

type Acceptor struct {
	mtx sync.Mutex
	//file to store the latest epoch
	epochFile *os.File
	//latest promise epoch
	latestEpoch uint64
	//file to store the accepted proposal
	proposalFile *os.File
	proposal     Proposal
}

type Proposal struct {
	Epoch uint64
	Value string
}

func NewAccepter(dir string) (*Acceptor, error) {
	epochFile, err := os.OpenFile(path.Join(dir, epochFileName), os.O_RDWR|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("open epoch file failed, dir[%s], error[%v]", dir, err))
	}

	proposalFile, err := os.OpenFile(path.Join(dir, proposalFileName), os.O_RDWR|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("open proposal file failed, dir[%s], error[%v]", dir, err))
	}

	acceptor := &Acceptor{
		epochFile:    epochFile,
		proposalFile: proposalFile,
	}

	if err := acceptor.recover(); err != nil {
		return nil, errors.New(fmt.Sprintf("recover acceptor state failed, error[%v]", err))
	}
	return acceptor, nil
}

func (acceptor *Acceptor) recover() error {
	//recover latest epoch
	epochRaw, err := ioutil.ReadAll(bufio.NewReaderSize(acceptor.epochFile, fileReadBufferSize))
	if err != nil {
		return err
	}
	if len(epochRaw) > 0 {
		epoch, err := strconv.ParseInt(string(epochRaw), 10, 64)
		if err != nil {
			return err
		}
		acceptor.latestEpoch = uint64(epoch)
	}

	//recover proposal
	proposalRaw, err := ioutil.ReadAll(bufio.NewReaderSize(acceptor.proposalFile, fileReadBufferSize))
	if err != nil {
		return err
	}
	if len(proposalRaw) > 0 {
		proposal := Proposal{}
		if err := json.Unmarshal(proposalRaw, &proposal); err != nil {
			return err
		}
		acceptor.proposal = proposal
	}
	return nil
}

/**
If the pEpoch is greater than any prepare request to which it has already responded,
then it responds to the request with a promise not to accept any more proposals whose
epoch less than pEpoch and with the highest-epoch proposal(if any) that it has accepted.
*/
func (acceptor *Acceptor) Prepare(pEpoch uint64) (bool, Proposal) {
	acceptor.mtx.Lock()
	defer acceptor.mtx.Unlock()

	if pEpoch <= acceptor.latestEpoch {
		return false, Proposal{}
	}

	if err := truncateAndWrite(acceptor.epochFile, fmt.Sprintf("%d", pEpoch)); err != nil {
		fmt.Println(fmt.Sprintf("truncateAndWrite file failed: error[%v]", err))
		return false, Proposal{}
	}
	acceptor.latestEpoch = pEpoch
	return true, acceptor.proposal
}

/**
this acceptor accept this proposal unless it has already responded to a prepare request
having a epoch greater than proposal.Epoch
*/
func (acceptor *Acceptor) Accept(proposal Proposal) bool {
	acceptor.mtx.Lock()
	defer acceptor.mtx.Unlock()

	if proposal.Epoch != acceptor.latestEpoch {
		return false
	}

	pStr, err := json.Marshal(proposal)
	if err != nil {
		fmt.Println(fmt.Sprintf("invalid proposal, json marshal fialed: error[%v]", err))
		return false
	}
	if err := truncateAndWrite(acceptor.proposalFile, string(pStr)); err != nil {
		fmt.Println(fmt.Sprintf("truncateAndWrite file failed: error[%v]", err))
		return false
	}
	acceptor.proposal = proposal
	return true
}

func truncateAndWrite(file *os.File, v string) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	written, err := file.WriteString(v)
	if err != nil {
		return err
	}
	if written != len(v) {
		return errors.New(fmt.Sprintf("file write failed: written[%d] != len[%d]", written, len(v)))
	}

	return nil
}
