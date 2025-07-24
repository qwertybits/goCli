package sjcli

type CLIData struct {
	arguments []string
	flags     map[string]any
}

func (data *CLIData) Argument(pos int) (string, bool) {
	if pos >= len(data.arguments) || pos < 0 {
		return "", false
	}
	return data.arguments[pos], true
}

func (data *CLIData) GetArguments() []string {
	return data.arguments
}

func (data *CLIData) GetString(flag string) (string, bool) {

	fl, exist := data.flags[flag]
	if !exist {
		return "", false
	}
	val, ok := fl.(string)
	if !ok {
		return "", false
	}
	return val, true
}

func (data *CLIData) GetBool(flag string) (bool, bool) {

	fl, exist := data.flags[flag]
	if !exist {
		return false, false
	}
	val, ok := fl.(bool)
	if !ok {
		return false, false
	}
	return val, true
}

func (data *CLIData) GetInt(flag string) (int, bool) {
	fl, exist := data.flags[flag]
	if !exist {
		return 0, false
	}
	val, ok := fl.(int)
	if !ok {
		return 0, false
	}
	return val, true
}
