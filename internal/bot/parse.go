package bot

import "strings"

func ParseArguments(input string) map[string]string {
	args := make(map[string]string)
	for _, arg := range strings.Split(input, "--") {
		arg = strings.TrimSpace(arg)

		if arg == "" {
			continue
		}

		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		args[parts[0]] = parts[1]
	}
	return args
}
