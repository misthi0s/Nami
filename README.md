<head>
<h1 align=center>Nami - Command & Control (C2)</h1>
</head>

<p align="center">
  <img src="images/nami.gif" alt="Nami"/>
</p>

Nami is a Command & Control (C2) framework, designed to help adversary emulation or red teamers with post-exploitation techniques. This is a very, very early release, with more features coming down the pipeline.

---

## Features

* HTTP C2 communication with encrypted and obfuscated payloads
* Built-in commands to perform common post-exploitation techniques, such as UAC bypassing or establishing persistence
* Cross-platform server infrastructure allowing it to run on Linux or Windows (payloads are currently Windows only)
* Support for multiple concurrent implants with the capability to switch between them at will
---
## Installation

Clone the repository:<br>
```git clone https://github.com/misthi0s/Nami```

Change the working directory:<br>
```cd Nami```

Build the project with Go:<br>
```go build -o Nami .```

Spin up the C2 server with custom IP address and port number (omitting these flags will default to 0.0.0.0:443):<br>
```./Nami -ip <IP_ADDR> -port <PORT>```

---
## Usage

Nami currently supports the following commands and flags:

<h3>Root</h3>
<table>
<tr>
<th>Command</th>
<th>Description</th>
</tr>
<tr>
<td>generate</td>
<td>Generate an implant payload</td>
</tr>
<tr>
<td>sessions</td>
<td>Interact with active Nami sessions</td>
</tr>
<tr>
<td>clear</td>
<td>Clear the screen output</td>
</tr>
<tr>
<td>help</td>
<td>Display help information; can be used against a specific command</td>
</tr>
<tr>
<td>exit</td>
<td>Exit Nami</td>
</tr>
<tr>
<td>help</td>
<td>Display help information; can be used against a specific command</td>
</tr>
</table>

<h3>Generate</h3>
<table>
<tr>
<th>Command</th>
<th>Description</th>
</tr>
<tr>
<td>-a/--arch</td>
<td>Choose architecture of implant (32 or 64) (Default: 32)</td>
</tr>
<tr>
<td>-d/--debug</td>
<td>Implement anti-debug techniques within the implant (boolean)</td>
</tr>
<tr>
<td>-i/--ip</td>
<td>IP address for implant to connect to</td>
</tr>
<tr>
<td>-p/--port</td>
<td>Port number for implant to connect to</td>
</tr>
<tr>
<td>-n/--name</td>
<td>Name of files or Registry keys written to system by certain commands (Default: Nami)</td>
</tr>
</table>

<h3>Sessions</h3>
<table>
<tr>
<th>Command</th>
<th>Description</th>
</tr>
<tr>
<td>-i/--interact</td>
<td>Choose session to interact with (running "sessions" will display all sessions)</td>
</tr>
</table>

<h3>Session Interaction</h3>
<table>
<tr>
<th>Command</th>
<th>Description</th>
</tr>
<tr>
<td>kill</td>
<td>Stop the implant process on the remote system</td>
</tr>
<tr>
<td>shell</td>
<td>Activate a cmd.exe shell on the remote system</td>
</tr>
<tr>
<td>persistence</td>
<td>Establish implant persistence on the remote system via a specific method</td>
</tr>
<tr>
<td>uacbypass</td>
<td>Establish a UAC bypass on the remote system via a specific method</td>
</tr>
</table>

<h3>Session - Persistence</h3>
<table>
<tr>
<th>Command</th>
<th>Description</th>
</tr>
<tr>
<td>--load</td>
<td>Create a Registry value for the legacy "Windows Load" persistence technique</td>
</tr>
<tr>
<td>--run</td>
<td>Create a Registry value for the "Run" persistence technique</td>
</tr>
<tr>
<td>--startupfolder</td>
<td>Create an LNK file in the system's StartUp directory for persistence</td>
</tr>
</table>

<h3>Session - UAC Bypass</h3>
<table>
<tr>
<th>Command</th>
<th>Description</th>
</tr>
<tr>
<td>--registry</td>
<td>Modify "EnableLUA" and associated "Prompt" Registry keys for persistence (requires admin)</td>
</tr>
</table>

---
## TODO

List of features that will be implemented in the future:

* Linux implant functionality
* Additional communication protocols (HTTPS, DNS, etc)
* Shellcode injection into another process
* File download/upload functionality
* Increased verbosity of the shell command, particularly with current directory
* More effective/efficient shell command executions
* Additional security bypasses
* Built-in Active Directory reconnaissance techniques

More TODOs will be added as development continues. If you have anything you'd like to see, feel free to open an issue or create a pull request with the desired functionality.

---
## Additional Notes

This is a very early alpha release of Nami. Due to this, there will be bugs and less features than other frameworks. Feel free to open any pull requests with bug issues or additional functionality.

---
## Issues

If you run into any issues with Nami, feel free to open an issue. 