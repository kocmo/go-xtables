package ebtables

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/kocmo/go-xtables"
	"github.com/kocmo/go-xtables/internal/xutil"
)

type onTableLine func(line []byte) (TableType, error)
type onChainLine func(line []byte) (*Chain, error)
type onRuleLine func(rule []byte, chain *Chain) (*Rule, error)

func (ebtables *EBTables) parse(data []byte, onTableLine onTableLine,
	onChainLine onChainLine, onRuleLine onRuleLine) (
	[]*Chain, []*Rule, error) {

	chains := []*Chain{}
	rules := []*Rule{}

	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)

	var tableType TableType
	var chain *Chain
	var err error
	var index int // index in current chain

	for scanner.Scan() {
		line := scanner.Bytes()
		if index == 0 {
			if bytes.HasPrefix(line, []byte("Bridge table")) {
				if onTableLine == nil {
					continue
				}
				tableType, err = onTableLine(line)
				if err != nil {
					ebtables.log.Errorf("parse table line err: %s", err)
					return nil, nil, err
				}
			} else if bytes.HasPrefix(line, []byte("Bridge chain")) {
				if onChainLine == nil {
					continue
				}
				chain, err = onChainLine(line)
				if err != nil {
					ebtables.log.Errorf("parse chain line err: %s", err)
					return nil, nil, err
				}
				chain.tableType = tableType
				chains = append(chains, chain)
			}
		} else {
			// rule or EOC(end of chain)
			if len(line) == 0 {
				index = 0
				continue
			}
			if onRuleLine == nil {
				continue
			}
			rule, err := onRuleLine(line, chain)
			if err != nil {
				ebtables.log.Errorf("parse rule line err: %s", err)
				return nil, nil, err
			}
			rule.tableType = tableType
			rules = append(rules, rule)
		}
		index++
	}
	return chains, rules, nil
}

func (ebtables *EBTables) parseTable(line []byte) (TableType, error) {
	buf := bytes.NewBuffer(line)
	_, err := buf.ReadString(':')
	if err != nil {
		return TableTypeNull, err
	}

	table := bytes.TrimFunc(buf.Bytes(), func(r rune) bool {
		if r == ' ' || r == '\n' {
			return true
		}
		return false
	})

	switch string(table) {
	case "nat":
		return TableTypeNat, nil
	case "filter":
		return TableTypeFilter, nil
	case "broute":
		return TableTypeBRoute, nil
	}
	return TableTypeNull, nil
}

func (ebtables *EBTables) parseChain(line []byte) (*Chain, error) {
	chain := &Chain{}
	buf := bytes.NewBuffer(line)
	_, err := buf.ReadString(':')
	if err != nil {
		ebtables.log.Errorf("parse chain read to first space err: %s", err)
		return nil, err
	}

	chain.chainType.name, err = buf.ReadString(',')
	if err != nil {
		ebtables.log.Errorf("parse chain read to second space err: %s", err)
		return nil, err
	}
	chain.chainType.name = strings.TrimSpace(chain.chainType.name[:len(chain.chainType.name)-1])
	switch chain.chainType.name {
	case "INPUT":
		chain.chainType = ChainTypeINPUT
	case "FORWARD":
		chain.chainType = ChainTypeFORWARD
	case "OUTPUT":
		chain.chainType = ChainTypeOUTPUT
	case "PREROUTING":
		chain.chainType = ChainTypePREROUTING
	case "BROUTING":
		chain.chainType = ChainTypeBROUTING
	case "POSTROUTING":
		chain.chainType = ChainTypePOSTROUTING
	default:
		userDefined := ChainTypeUserDefined
		userDefined.name = chain.chainType.name
		chain.chainType = userDefined
	}

	rest := buf.Bytes()
	attrs := bytes.FieldsFunc(rest, func(r rune) bool {
		if r == ',' || r == ' ' {
			return true
		}
		return false
	})
	if len(attrs)%2 != 0 {
		return nil, xtables.ErrChainAttrsNotRecognized
	}

	chain.policy = newTargetAccept()

	pairs := len(attrs) / 2
	for i := 0; i < pairs; i++ {
		index := i * 2
		first := attrs[index]
		second := attrs[index+1]

		// entries
		if bytes.HasPrefix(bytes.TrimSpace(first), []byte("entries")) {
			num, err := strconv.Atoi(string(second))
			if err != nil {
				return nil, err
			}
			chain.entries = num
		}

		// policy
		if bytes.HasPrefix(bytes.TrimSpace(first), []byte("policy")) {
			switch string(second) {
			case "ACCEPT":
			case "DROP":
				chain.policy = newTargetDrop()
			case "RETURN":
				chain.policy = newTargetReturn()
			}
		}
	}
	return chain, nil
}

func (ebtables *EBTables) parseRule(line []byte, chain *Chain) (*Rule, error) {
	rule := &Rule{
		chain:      chain,
		matchMap:   map[MatchType]Match{},
		optionMap:  map[OptionType]Option{},
		watcherMap: map[WatcherType]Watcher{},
		lineNumber: -1,
	}
	delimiter := []byte{'.', ' '}
	index := bytes.Index(line, delimiter)
	if index > 0 && index < len(line) {
		ln, err := strconv.Atoi(strings.TrimSpace(string(line[:index])))
		if err == nil {
			rule.lineNumber = ln
			line = line[index+len(delimiter):]
		}
	}
	// then matches
	matches, index, err := ebtables.parseMatch(line)
	if err != nil {
		ebtables.log.Errorf("parse match: %s err: %s", string(line), err)
		return nil, err
	}
	for _, match := range matches {
		rule.matchMap[match.Type()] = match
	}
	line = line[index:]

	// watcher
	watchers, index, err := ebtables.parseWatcher(line)
	if err != nil {
		ebtables.log.Errorf("parse watcher: %s err: %s", string(line), err)
		return nil, err
	}
	for _, watcher := range watchers {
		rule.watcherMap[watcher.Type()] = watcher
	}
	line = line[index:]

	// then target
	target, index, err := ebtables.parseTarget(line)
	if err != nil {
		if err != xtables.ErrTargetNotFound {
			ebtables.log.Errorf("parse target: %s err: %s", string(line), err)
			return nil, err
		}
		rule.target = nil
	} else {
		rule.target = target
		line = line[index:]
	}

	// then pkt and bytes count
	pcnt, bcnt, ok := parsePktsAndBytes(line)
	if ok {
		opt, _ := newOptionCounters(pcnt, bcnt)
		rule.optionMap[opt.Type()] = opt
		rule.packetCounter = pcnt
		rule.byteCounter = bcnt
	}

	return rule, nil
}

func parsePktsAndBytes(params []byte) (int64, int64, bool) {
	pattern := `,? *pcnt = ([0-9A-Za-z]+) -- bcnt = ([0-9A-Za-z]+) *`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindSubmatch(params)
	if len(matches) != 3 {
		return 0, 0, false
	}
	pcnt, err := xutil.UnfoldDecimal(string(matches[1]))
	if err != nil {
		return 0, 0, false
	}
	bcnt, err := xutil.UnfoldDecimal(string(matches[2]))
	if err != nil {
		return 0, 0, false
	}
	return pcnt, bcnt, true
}
