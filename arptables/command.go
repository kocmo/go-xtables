package arptables

import "strconv"

type CommandType int

func (ct CommandType) Type() string {
	return "CommandType"
}

func (ct CommandType) Value() string {
	return strconv.Itoa(int(ct))
}

const (
	_                      CommandType = iota
	CommandTypeAppend                  // append
	CommandTypeDelete                  // delete
	CommandTypeInsert                  // insert
	CommandTypeFlush                   // flush
	CommandTypePolicy                  // policy
	CommandTypeZero                    // zero
	CommandTypeList                    // list
	CommandTypeNewChain                // new_chain
	CommandTypeDeleteChain             // delete_chain
	CommandTypeRenameChain             // rename_chain
)

type HasRulenum interface {
	Rulenum() uint32
}

type Command interface {
	Type() CommandType
	Short() string
	Long() string
}

type baseCommand struct {
	commandType CommandType
	child       Command
}

func (bc baseCommand) setChild(child Command) {
	bc.child = child
}

func (bc baseCommand) Type() CommandType {
	return bc.commandType
}

func (bc baseCommand) Short() string {
	if bc.child != nil {
		return bc.child.Short()
	}
	return ""
}

func (bc baseCommand) Long() string {
	if bc.child != nil {
		return bc.child.Long()
	}
	return bc.Short()
}

func NewAppend() *Append {
	command := &Append{
		baseCommand: baseCommand{
			commandType: CommandTypeAppend,
		},
	}
	command.setChild(command)
	return command
}

type Append struct {
	baseCommand
}

func (apd *Append) Short() string {
	return "-A"
}

func (apd *Append) Long() string {
	return "--append"
}

func NewDelete(rulenum uint32) *Delete {
	command := &Delete{
		baseCommand: baseCommand{
			commandType: CommandTypeDelete,
		},
		rnum: rulenum,
	}
	command.setChild(command)
	return command
}

type Delete struct {
	baseCommand
	rnum uint32
}

func (del *Delete) Rulenum() uint32 {
	return del.rnum
}

func (del *Delete) Short() string {
	return "-D"
}

func (del *Delete) Long() string {
	return "--delete"
}

func NewInsert(rulenum uint32) *Insert {
	command := &Insert{
		baseCommand: baseCommand{
			commandType: CommandTypeInsert,
		},
		rnum: rulenum,
	}
	command.setChild(command)
	return command
}

type Insert struct {
	baseCommand
	rnum uint32
}

func (insert *Insert) Rulenum() uint32 {
	return insert.rnum
}

func (insert *Insert) Short() string {
	return "-I"
}

func (insert *Insert) Long() string {
	return "--insert"
}

func NewList() *List {
	command := &List{
		baseCommand: baseCommand{
			commandType: CommandTypeList,
		},
	}
	command.setChild(command)
	return command
}

type List struct {
	baseCommand
}

func (list *List) Short() string {
	return "-L"
}

func (list *List) Long() string {
	return "--list"
}

func NewFlush() *Flush {
	command := &Flush{
		baseCommand: baseCommand{
			commandType: CommandTypeFlush,
		},
	}
	command.setChild(command)
	return command
}

type Flush struct {
	baseCommand
}

func (flush *Flush) Short() string {
	return "-F"
}

func (flush *Flush) Long() string {
	return "--flush"
}

func NewPolicy() *Policy {
	command := &Policy{
		baseCommand: baseCommand{
			commandType: CommandTypePolicy,
		},
	}
	command.setChild(command)
	return command
}

type Policy struct {
	baseCommand
	targetType TargetType
}

func (policy *Policy) Short() string {
	return "-P"
}

func (policy *Policy) Long() string {
	return "--policy"
}

func NewZero() *Zero {
	command := &Zero{
		baseCommand: baseCommand{
			commandType: CommandTypeZero,
		},
	}
	command.setChild(command)
	return command
}

type Zero struct {
	baseCommand
	rnum uint32
}

func (zero *Zero) Rulenum() uint32 {
	return zero.rnum
}

func (zero *Zero) Short() string {
	return "-Z"
}

func (zero *Zero) Long() string {
	return "--zero"
}

func NewNewChain() *NewChain {
	command := &NewChain{
		baseCommand: baseCommand{
			commandType: CommandTypeNewChain,
		},
	}
	command.setChild(command)
	return command
}

type NewChain struct {
	baseCommand
}

func (nc *NewChain) Short() string {
	return "-N"
}

func (nc *NewChain) Long() string {
	return "--new-chain"
}

func NewDeleteChain() *DeleteChain {
	command := &DeleteChain{
		baseCommand: baseCommand{
			commandType: CommandTypeDeleteChain,
		},
	}
	command.setChild(command)
	return command
}

type DeleteChain struct {
	baseCommand
}

func (dc *DeleteChain) Short() string {
	return "-X"
}

func (dc *DeleteChain) Long() string {
	return "--delete-chain"
}

func NewRenameChain() *RenameChain {
	command := &RenameChain{
		baseCommand: baseCommand{
			commandType: CommandTypeRenameChain,
		},
	}
	command.setChild(command)
	return command
}

type RenameChain struct {
	baseCommand
	newChain string // user supplied name.
}

func (renameChain *RenameChain) Short() string {
	return "-E"
}

func (renameChain *RenameChain) Long() string {
	return "--rename-chain"
}
