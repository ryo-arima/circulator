package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// MCode represents a message code with predefined messages
// Follows the pattern from https://github.com/ryo-arima/vem/blob/main/src/util/mcode.rs
type MCode struct {
	Code    string
	Message string
}

// FormatWithOptional formats the message with optional additional message
func (m MCode) FormatWithOptional(optionalMessage string) string {
	if optionalMessage == "" {
		return m.Message
	}
	return fmt.Sprintf("%s: %s", m.Message, optionalMessage)
}

// Configuration Layer Codes - C_* (Config)
var (
	CL1 = MCode{"C-L1", "Configuration load success"}
	CL2 = MCode{"C-L2", "Configuration load failed"}
	CL3 = MCode{"C-L3", "Configuration parse error"}
	CS1 = MCode{"C-S1", "Configuration save success"}
	CS2 = MCode{"C-S2", "Configuration save failed"}
	CV1 = MCode{"C-V1", "Configuration validation success"}
	CV2 = MCode{"C-V2", "Configuration validation failed"}
)

// Client Repository Layer Codes - CR_* (Client Repository)
var (
	CRI1 = MCode{"CR-I1", "Client repository initialization success"}
	CRI2 = MCode{"CR-I2", "Client repository initialization failed"}
	CRS1 = MCode{"CR-S1", "Client-Server HTTP communication success"}
	CRS2 = MCode{"CR-S2", "Client-Server HTTP communication failed"}
	CRA1 = MCode{"CR-A1", "Client-Agent HTTP communication success"}
	CRA2 = MCode{"CR-A2", "Client-Agent HTTP communication failed"}
	CRP1 = MCode{"CR-P1", "Client-Pulsar messaging success"}
	CRP2 = MCode{"CR-P2", "Client-Pulsar messaging failed"}
)

// Client Controller Layer Codes - CC_* (Client Controller)
var (
	CCI1 = MCode{"CC-I1", "Client controller initialization success"}
	CCI2 = MCode{"CC-I2", "Client controller initialization failed"}
	CCR1 = MCode{"CC-R1", "Client request processing success"}
	CCR2 = MCode{"CC-R2", "Client request processing failed"}
)

// Agent Repository Layer Codes - AR_* (Agent Repository)
var (
	ARI1 = MCode{"AR-I1", "Agent repository initialization success"}
	ARI2 = MCode{"AR-I2", "Agent repository initialization failed"}
	ARS1 = MCode{"AR-S1", "Agent-Server HTTP communication success"}
	ARS2 = MCode{"AR-S2", "Agent-Server HTTP communication failed"}
	ARP1 = MCode{"AR-P1", "Agent Pulsar producer success"}
	ARP2 = MCode{"AR-P2", "Agent Pulsar producer failed"}
	ARP3 = MCode{"AR-P3", "Agent Pulsar consumer success"}
	ARP4 = MCode{"AR-P4", "Agent Pulsar consumer failed"}
	ARL1 = MCode{"AR-L1", "Agent local system access success"}
	ARL2 = MCode{"AR-L2", "Agent local system access failed"}
)

// Agent Controller Layer Codes - AC_* (Agent Controller)
var (
	ACI1 = MCode{"AC-I1", "Agent controller initialization success"}
	ACI2 = MCode{"AC-I2", "Agent controller initialization failed"}
	ACS1 = MCode{"AC-S1", "Agent-Server gRPC handling success"}
	ACS2 = MCode{"AC-S2", "Agent-Server gRPC handling failed"}
	ACC1 = MCode{"AC-C1", "Agent-Client HTTP handling success"}
	ACC2 = MCode{"AC-C2", "Agent-Client HTTP handling failed"}
	ACP1 = MCode{"AC-P1", "Agent Pulsar message handling success"}
	ACP2 = MCode{"AC-P2", "Agent Pulsar message handling failed"}
)

// Server Repository Layer Codes - SR_* (Server Repository)
var (
	SRI1 = MCode{"SR-I1", "Server repository initialization success"}
	SRI2 = MCode{"SR-I2", "Server repository initialization failed"}
	SRM1 = MCode{"SR-M1", "Server MySQL operation success"}
	SRM2 = MCode{"SR-M2", "Server MySQL operation failed"}
	SRP1 = MCode{"SR-P1", "Server Pulsar messaging success"}
	SRP2 = MCode{"SR-P2", "Server Pulsar messaging failed"}
	SRA1 = MCode{"SR-A1", "Server-Agent gRPC communication success"}
	SRA2 = MCode{"SR-A2", "Server-Agent gRPC communication failed"}
)

// Server Controller Layer Codes - SC_* (Server Controller)
var (
	SCI1 = MCode{"SC-I1", "Server controller initialization success"}
	SCI2 = MCode{"SC-I2", "Server controller initialization failed"}
	SCC1 = MCode{"SC-C1", "Server-Client HTTP handling success"}
	SCC2 = MCode{"SC-C2", "Server-Client HTTP handling failed"}
	SCA1 = MCode{"SC-A1", "Server-Agent HTTP handling success"}
	SCA2 = MCode{"SC-A2", "Server-Agent HTTP handling failed"}
	SCP1 = MCode{"SC-P1", "Server Pulsar handling success"}
	SCP2 = MCode{"SC-P2", "Server Pulsar handling failed"}
)

// System Layer Codes - SY_* (SYstem)
var (
	SYS1 = MCode{"SY-S1", "Application started"}
	SYS2 = MCode{"SY-S2", "Application terminated successfully"}
	SYS3 = MCode{"SY-S3", "Application terminated with error"}
	SYE1 = MCode{"SY-E1", "Unexpected error occurred"}
)

// File System Operation Codes - FS_* (File System)
var (
	FSO1 = MCode{"FS-O1", "File system operation success"}
	FSO2 = MCode{"FS-O2", "File system operation failed"}
	FSM1 = MCode{"FS-M1", "Directory creation success"}
	FSM2 = MCode{"FS-M2", "Directory creation failed"}
	FSW1 = MCode{"FS-W1", "File write success"}
	FSW2 = MCode{"FS-W2", "File write failed"}
	FSR1 = MCode{"FS-R1", "File read success"}
	FSR2 = MCode{"FS-R2", "File read failed"}
)

// Communication Pattern Codes - CP_* (Communication Pattern)
var (
	CP01 = MCode{"CP-01", "Direct client-server communication"}
	CP02 = MCode{"CP-02", "Direct client-agent communication"}
	CP03 = MCode{"CP-03", "Server-agent gRPC communication"}
	CP04 = MCode{"CP-04", "Agent-server HTTP communication"}
	CP05 = MCode{"CP-05", "Client-server-MySQL access"}
	CP06 = MCode{"CP-06", "Agent-server-MySQL access"}
	CP07 = MCode{"CP-07", "Direct server-MySQL access"}
	CP08 = MCode{"CP-08", "Client-Pulsar messaging"}
	CP09 = MCode{"CP-09", "Agent-Pulsar messaging"}
	CP10 = MCode{"CP-10", "Server-Pulsar messaging"}
	CP11 = MCode{"CP-11", "Pulsar-client messaging"}
	CP12 = MCode{"CP-12", "Pulsar-agent messaging"}
	CP13 = MCode{"CP-13", "Pulsar-server messaging"}
	CP14 = MCode{"CP-14", "Client-Pulsar-server communication"}
	CP15 = MCode{"CP-15", "Client-Pulsar-agent communication"}
	CP16 = MCode{"CP-16", "Agent-Pulsar-client communication"}
	CP17 = MCode{"CP-17", "Agent-Pulsar-server communication"}
	CP18 = MCode{"CP-18", "Server-Pulsar-client communication"}
	CP19 = MCode{"CP-19", "Server-Pulsar-agent communication"}
	CP20 = MCode{"CP-20", "Client-server-Pulsar-agent flow"}
	CP21 = MCode{"CP-21", "Agent-server-Pulsar-client flow"}
	CP22 = MCode{"CP-22", "External-agent-Pulsar-server flow"}
	CP23 = MCode{"CP-23", "External-agent-server-Pulsar-client flow"}
)

// Additional Legacy Code Mappings - for backward compatibility
var (
	// Server Router codes
	SRIR   = MCode{"SR-IR", "Initializing HTTP router"}
	SRCARI = MCode{"SR-CARI", "Controllers and repositories initialized"}
	SRRAE  = MCode{"SR-RAE", "Registering authentication endpoints"}
	SRRPAE = MCode{"SR-RPAE", "Registering protected API endpoints"}
	SRHRIS = MCode{"SR-HRIS", "HTTP router initialized successfully"}

	// Agent Controller Common codes
	ACCGS = MCode{"ACC-GS", "Status requested via controller"}
	ACCP  = MCode{"ACC-P", "Ping requested"}

	// Agent Repository Pulsar Consumer codes
	APCERM  = MCode{"APC-RM", "Error receiving message"}
	APCPSD  = MCode{"APC-PSD", "Error parsing stream data"}
	APCPSD2 = MCode{"APC-PSD2", "Error processing stream data"}

	// Server Middleware codes
	SML = MCode{"SM-L", "HTTP request processed"}

	// Agent Register codes
	ARSGSR    = MCode{"AR-SGSR", "Starting gRPC service registration"}
	ARIC      = MCode{"AR-IC", "Initializing controllers"}
	ARFISC    = MCode{"AR-FISC", "Failed to initialize stream controller"}
	ARGRPC    = MCode{"AR-GRPC", "gRPC services registration setup completed"}
	ARSGRPC   = MCode{"AR-SGRPC", "Starting gRPC server"}
	ARFTLOP   = MCode{"AR-FTLOP", "Failed to listen on port"}
	ARGRPCS   = MCode{"AR-GRPCS", "gRPC server starting"}
	ARFTSGRPC = MCode{"AR-FTSGRPC", "Failed to serve gRPC server"}

	// Client Base codes
	CBCE  = MCode{"CB-CE", "Command executed"}
	CBIBC = MCode{"CB-IBC", "Initializing base commands"}
	CBBCC = MCode{"CB-BCC", "Bootstrap command called"}
	CBCCC = MCode{"CB-CCC", "Create command called"}
	CBSCA = MCode{"CB-SCA", "Starting client application"}
	CBACR = MCode{"CB-ACR", "All commands registered"}

	// Client Controller codes
	CCLCE1  = MCode{"CC-LC-E1", "Failed to get email flag"}
	CCLCE2  = MCode{"CC-LC-E2", "Failed to get password flag"}
	CCRTCE1 = MCode{"CC-RTC-E1", "Failed to get refresh-token flag"}
	CCLOCE1 = MCode{"CC-LOC-E1", "Failed to get access-token flag"}
	CCVTCE1 = MCode{"CC-VTC-E1", "Failed to get access-token flag"}
	CCUICE1 = MCode{"CC-UIC-E1", "Failed to get access-token flag"}

	// Agent UseCase codes
	AUCHC  = MCode{"AUC-HC", "Performing agent health check"}
	AUASPC = MCode{"AUA-SPC", "Setting processing config"}
	AUAGPC = MCode{"AUA-GPC", "Getting processing config"}
	AUAPAD = MCode{"AUA-PAD", "Processed agent data"}

	// Agent Controller Agent codes
	ACAPSD = MCode{"ACA-PSD", "Processing stream data via controller"}

	// Agent Repository API codes
	AREGA    = MCode{"ARA-GA", "Getting all agents"}
	ARECA    = MCode{"ARA-CA", "Counting agents"}
	ARECRA   = MCode{"ARA-CRA", "Creating agent"}
	AREUA    = MCode{"ARA-UA", "Updating agent"}
	AREDA    = MCode{"ARA-DA", "Deleting agent"}
	AREGAI   = MCode{"ARA-GAI", "Getting agent info"}
	ARECRAI  = MCode{"ARA-CRAI", "Creating agent info"}
	AREUAI   = MCode{"ARA-UAI", "Updating agent info"}
	AREDAI   = MCode{"ARA-DAI", "Deleting agent info"}
	AREGAS   = MCode{"ARA-GAS", "Getting agent system info"}
	ARECRAS  = MCode{"ARA-CRAS", "Creating agent system info"}
	AREUAS   = MCode{"ARA-UAS", "Updating agent system info"}
	AREDAS   = MCode{"ARA-DAS", "Deleting agent system info"}
	AREGAC   = MCode{"ARA-GAC", "Getting agent config"}
	ARECRAC  = MCode{"ARA-CRAC", "Creating agent config"}
	AREUAC   = MCode{"ARA-UAC", "Updating agent config"}
	AREDAC   = MCode{"ARA-DAC", "Deleting agent config"}
	AREGACR  = MCode{"ARA-GACR", "Getting agent config rules"}
	ARECRACR = MCode{"ARA-CRACR", "Creating agent config rules"}
	AREUACR  = MCode{"ARA-UACR", "Updating agent config rules"}
	AREDACR  = MCode{"ARA-DACR", "Deleting agent config rules"}
	ARESPC   = MCode{"ARA-SPC", "Setting processing config"}
	AREGPC   = MCode{"ARA-GPC", "Getting processing config"}
)

// Agent Repository API Common codes
var (
	ARACLOG = MCode{"ARAC-LOG", "Authentication operation"}
	ARACREG = MCode{"ARAC-REG", "Agent registration operation"}
	ARACHB  = MCode{"ARAC-HB", "Heartbeat operation"}
	ARACRT  = MCode{"ARAC-RT", "Token refresh operation"}
)

// Server UseCase Agent codes
var (
	SUAGA   = MCode{"SUA-GA", "Getting all agents"}
	SUACA   = MCode{"SUA-CA", "Counting agents"}
	SUAGA1  = MCode{"SUA-GA1", "Getting agent by UUID"}
	SUACA1  = MCode{"SUA-CA1", "Creating new agent"}
	SUAUA   = MCode{"SUA-UA", "Updating agent"}
	SUADA   = MCode{"SUA-DA", "Deleting agent"}
	SUAGAI  = MCode{"SUA-GAI", "Getting agent info"}
	SUACAI  = MCode{"SUA-CAI", "Creating agent info"}
	SUAUAI  = MCode{"SUA-UAI", "Updating agent info"}
	SUADAI  = MCode{"SUA-DAI", "Deleting agent info"}
	SUAGAS  = MCode{"SUA-GAS", "Getting agent system info"}
	SUACAS  = MCode{"SUA-CAS", "Creating agent system info"}
	SUAUAS  = MCode{"SUA-UAS", "Updating agent system info"}
	SUADAS  = MCode{"SUA-DAS", "Deleting agent system info"}
	SUAGSPC = MCode{"SUA-GSPC", "Getting stream processing config"}
	SUACSPC = MCode{"SUA-CSPC", "Creating stream processing config"}
	SUAUSPC = MCode{"SUA-USPC", "Updating stream processing config"}
	SUADSPC = MCode{"SUA-DSPC", "Deleting stream processing config"}
	SUAGPR  = MCode{"SUA-GPR", "Getting processing rules"}
	SUACPR  = MCode{"SUA-CPR", "Creating processing rule"}
	SUAUPR  = MCode{"SUA-UPR", "Updating processing rule"}
	SUADPR  = MCode{"SUA-DPR", "Deleting processing rule"}
)

// Server UseCase Common codes
var (
	SUCVU = MCode{"SUC-VU", "Validating user credentials"}
	SUCGT = MCode{"SUC-GT", "Generating token for user"}
	SUCVT = MCode{"SUC-VT", "Validating token"}
	SUCLU = MCode{"SUC-LU", "Processing user login"}
	SUCRU = MCode{"SUC-RU", "Processing token refresh"}
)

// Agent Base codes
var (
	ABM    = MCode{"AB-M", "Starting Agent"}
	ABME2  = MCode{"AB-M-E2", "Failed to register agent"}
	ABME3  = MCode{"AB-M-E3", "Failed to start gRPC server"}
	ABRA   = MCode{"AB-RA", "Registering agent with server"}
	ABRAE3 = MCode{"AB-RA-E3", "Failed to get system info"}
	ABRAE4 = MCode{"AB-RA-E4", "Failed to register agent"}
	ABRAS  = MCode{"AB-RA-S", "Agent registration completed"}
)

// Client Repository Server codes
var (
	CRSLOGIN = MCode{"CRS-LOGIN", "Client server login attempt"}
	CRSGA    = MCode{"CRS-GA", "Client getting all agents from server"}
	CRSCA    = MCode{"CRS-CA", "Client creating agent on server"}
	CRSSUCC  = MCode{"CRS-SUCC", "Client server operation successful"}
	CRSERR   = MCode{"CRS-ERR", "Client server operation error"}
)

// Client Repository Pulsar codes
var (
	CRPINIT  = MCode{"CRP-INIT", "Client Pulsar repository initialized"}
	CRPPUB   = MCode{"CRP-PUB", "Client publishing to Pulsar"}
	CRPCONS  = MCode{"CRP-CONS", "Client consuming from Pulsar"}
	CRPREC   = MCode{"CRP-REC", "Client received Pulsar message"}
	CRPSTOP  = MCode{"CRP-STOP", "Client stopping Pulsar consumption"}
	CRPCLOSE = MCode{"CRP-CLOSE", "Client Pulsar repository closed"}
	CRPSUCC  = MCode{"CRP-SUCC", "Client Pulsar operation successful"}
	CRPERR   = MCode{"CRP-ERR", "Client Pulsar operation error"}
)

// Server Repository MySQL codes
var (
	SRMCONN     = MCode{"SRM-CONN", "Server MySQL connection"}
	SRMINIT     = MCode{"SRM-INIT", "Server MySQL repository initialized"}
	SRMMIG      = MCode{"SRM-MIG", "Server MySQL migrations"}
	SRMGA       = MCode{"SRM-GA", "Server MySQL get all agents"}
	SRMGAU      = MCode{"SRM-GAU", "Server MySQL get agent by UUID"}
	SRMCA       = MCode{"SRM-CA", "Server MySQL create agent"}
	SRMUA       = MCode{"SRM-UA", "Server MySQL update agent"}
	SRMDA       = MCode{"SRM-DA", "Server MySQL delete agent"}
	SRMCOUNT    = MCode{"SRM-COUNT", "Server MySQL count agents"}
	SRMCLOSE    = MCode{"SRM-CLOSE", "Server MySQL repository closed"}
	SRMNOTFOUND = MCode{"SRM-NOTFOUND", "Server MySQL record not found"}
	SRMSUCC     = MCode{"SRM-SUCC", "Server MySQL operation successful"}
	SRMERR      = MCode{"SRM-ERR", "Server MySQL operation error"}
)

// Server Repository Pulsar codes
var (
	SRPINIT  = MCode{"SRP-INIT", "Server Pulsar repository initialized"}
	SRPPUB   = MCode{"SRP-PUB", "Server publishing to Pulsar"}
	SRPNOT   = MCode{"SRP-NOT", "Server publishing notification"}
	SRPCONS  = MCode{"SRP-CONS", "Server consuming from Pulsar"}
	SRPREC   = MCode{"SRP-REC", "Server received Pulsar message"}
	SRPSTOP  = MCode{"SRP-STOP", "Server stopping Pulsar consumption"}
	SRPCLOSE = MCode{"SRP-CLOSE", "Server Pulsar repository closed"}
	SRPSUCC  = MCode{"SRP-SUCC", "Server Pulsar operation successful"}
	SRPERR   = MCode{"SRP-ERR", "Server Pulsar operation error"}
)

// Agent Repository API Server codes
var (
	ARSINIT  = MCode{"ARS-INIT", "Agent Server API repository initialized"}
	ARSLOGIN = MCode{"ARS-LOGIN", "Agent attempting server login"}
	ARSGINFO = MCode{"ARS-GINFO", "Agent getting info from server"}
	ARSSPREP = MCode{"ARS-SPREP", "Agent sending status report"}
	ARSGREG  = MCode{"ARS-GREG", "Agent getting registration from server"}
	ARSSREG  = MCode{"ARS-SREG", "Agent sending registration to server"}
	ARSSUCC  = MCode{"ARS-SUCC", "Agent Server API operation successful"}
	ARSERR   = MCode{"ARS-ERR", "Agent Server API operation error"}
)

// Agent Repository Pulsar Producer codes
var (
	ARPPINIT = MCode{"ARPP-INIT", "Agent Pulsar producer initialized"}
	ARPPUB   = MCode{"ARPP-PUB", "Agent publishing to Pulsar"}
	ARPREP   = MCode{"ARPP-REP", "Agent publishing report"}
	ARPNOT   = MCode{"ARPP-NOT", "Agent publishing notification"}
	ARPCLOSE = MCode{"ARPP-CLOSE", "Agent Pulsar producer closed"}
	ARPSUCC  = MCode{"ARPP-SUCC", "Agent Pulsar producer operation successful"}
	ARPERR   = MCode{"ARPP-ERR", "Agent Pulsar producer operation error"}
)

// Agent Repository Pulsar Consumer codes
var (
	ARCCINIT = MCode{"ARCC-INIT", "Agent Pulsar consumer initialized"}
	ARCCONS  = MCode{"ARCC-CONS", "Agent consuming from Pulsar"}
	ARCREC   = MCode{"ARCC-REC", "Agent received Pulsar message"}
	ARCPROC  = MCode{"ARCC-PROC", "Agent processing Pulsar message"}
	ARCSTOP  = MCode{"ARCC-STOP", "Agent stopping Pulsar consumption"}
	ARCCLOSE = MCode{"ARCC-CLOSE", "Agent Pulsar consumer closed"}
	ARCSUCC  = MCode{"ARCC-SUCC", "Agent Pulsar consumer operation successful"}
	ARCERR   = MCode{"ARCC-ERR", "Agent Pulsar consumer operation error"}
)

// Agent Repository Local System codes
var (
	ALSINIT  = MCode{"ALS-INIT", "Agent local system repository initialized"}
	ALSGINFO = MCode{"ALS-GINFO", "Agent getting system info"}
	ALSGSTAT = MCode{"ALS-GSTAT", "Agent getting system status"}
	ALSGREG  = MCode{"ALS-GREG", "Agent getting registration info"}
	ALSSREG  = MCode{"ALS-SREG", "Agent storing registration info"}
	ALSSUCC  = MCode{"ALS-SUCC", "Agent local system operation successful"}
	ALSERR   = MCode{"ALS-ERR", "Agent local system operation error"}
)

// LogLevel represents the log level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a structured log entry for fluentd compatibility
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Code      string                 `json:"code"`      // ログコード（例: PSR-IR）
	Component string                 `json:"component"` // client, server, agent
	Service   string                 `json:"service"`   // specific service name
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	File      string                 `json:"file,omitempty"`
	Function  string                 `json:"function,omitempty"`
	Line      int                    `json:"line,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	AgentID   string                 `json:"agent_id,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Component    string `json:"component" yaml:"component"`         // client, server, agent
	Service      string `json:"service" yaml:"service"`             // specific service name
	Level        string `json:"level" yaml:"level"`                 // DEBUG, INFO, WARN, ERROR, FATAL
	Structured   bool   `json:"structured" yaml:"structured"`       // JSON format for fluentd
	EnableCaller bool   `json:"enable_caller" yaml:"enable_caller"` // Include file/line info
	Output       string `json:"output" yaml:"output"`               // stdout, stderr, file path
}

// Logger represents the application logger with dependency injection capability
type Logger struct {
	config     *LoggerConfig
	level      LogLevel
	output     io.Writer
	baseConfig *BaseConfig // Dependency injection of BaseConfig
}

// LoggerInterface defines the logging interface for dependency injection
//
// Usage Examples:
//
// Using MCode (recommended):
//
//	logger.InfoWithMCode(SYS1, "", map[string]interface{}{"port": "8080"})
//	logger.ErrorWithMCode(SRM2, "connection timeout", map[string]interface{}{"host": "localhost"})
//	logger.InfoWithMCode(ARP1, "sent to topic: user-events", nil)
//
// Using legacy methods (backward compatibility):
//
//	logger.Info("SYS-S1", "Application started", map[string]interface{}{"port": "8080"})
type LoggerInterface interface {
	// MCode-based logging methods
	DEBUG(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	INFO(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	WARN(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	ERROR(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	FATAL(mcode MCode, optionalMessage string, fields ...map[string]interface{})
}

// NewLogger creates a new logger instance with BaseConfig dependency injection
func NewLogger(loggerConfig LoggerConfig, baseConfig *BaseConfig) LoggerInterface {
	logger := &Logger{
		config:     &loggerConfig,
		baseConfig: baseConfig,
		output:     os.Stdout,
	}

	// Set log level
	switch strings.ToUpper(loggerConfig.Level) {
	case "DEBUG":
		logger.level = DEBUG
	case "INFO":
		logger.level = INFO
	case "WARN":
		logger.level = WARN
	case "ERROR":
		logger.level = ERROR
	case "FATAL":
		logger.level = FATAL
	default:
		logger.level = INFO
	}

	// Set output
	switch loggerConfig.Output {
	case "stderr":
		logger.output = os.Stderr
	case "stdout", "":
		logger.output = os.Stdout
	default:
		// File output
		if file, err := os.OpenFile(loggerConfig.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666); err == nil {
			logger.output = file
		} else {
			logger.output = os.Stdout
			// Use MCode for error logging
			logger.ERROR(FSW2, fmt.Sprintf("file: %s, error: %s", loggerConfig.Output, err.Error()))
		}
	}

	return logger
}

// log writes a log entry using MCode
func (l *Logger) log(level LogLevel, mcode MCode, optionalMessage string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	finalMessage := mcode.FormatWithOptional(optionalMessage)

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level.String(),
		Code:      mcode.Code,
		Component: l.config.Component,
		Service:   l.config.Service,
		Message:   finalMessage,
		Fields:    fields,
	}

	l.writeLogEntry(entry)
}

// writeLogEntry writes the actual log entry to output
func (l *Logger) writeLogEntry(entry LogEntry) {
	// Add caller information if enabled
	if l.config.EnableCaller {
		if pc, file, line, ok := runtime.Caller(4); ok { // Adjusted for additional call stack
			entry.File = file
			entry.Line = line
			if fn := runtime.FuncForPC(pc); fn != nil {
				entry.Function = fn.Name()
			}
		}
	}

	// Extract common fields from fields map
	if entry.Fields != nil {
		if traceID, ok := entry.Fields["trace_id"].(string); ok {
			entry.TraceID = traceID
			delete(entry.Fields, "trace_id")
		}
		if requestID, ok := entry.Fields["request_id"].(string); ok {
			entry.RequestID = requestID
			delete(entry.Fields, "request_id")
		}
		if userID, ok := entry.Fields["user_id"].(string); ok {
			entry.UserID = userID
			delete(entry.Fields, "user_id")
		}
		if agentID, ok := entry.Fields["agent_id"].(string); ok {
			entry.AgentID = agentID
			delete(entry.Fields, "agent_id")
		}
		if err, ok := entry.Fields["error"].(string); ok {
			entry.Error = err
			delete(entry.Fields, "error")
		}
		if err, ok := entry.Fields["error"].(error); ok {
			entry.Error = err.Error()
			delete(entry.Fields, "error")
		}
	}

	if l.config.Structured {
		// JSON format for fluentd
		if jsonBytes, err := json.Marshal(entry); err == nil {
			fmt.Fprintln(l.output, string(jsonBytes))
		} else {
			// Fallback to simple format
			fmt.Fprintf(l.output, "[%s] %s [%s] %s/%s: %s\n",
				entry.Timestamp, entry.Level, entry.Code, entry.Component, entry.Service, entry.Message)
		}
	} else {
		// Human-readable format
		fmt.Fprintf(l.output, "[%s] %s [%s] %s/%s: %s",
			entry.Timestamp, entry.Level, entry.Code, entry.Component, entry.Service, entry.Message)
		if len(entry.Fields) > 0 {
			if fieldsJSON, err := json.Marshal(entry.Fields); err == nil {
				fmt.Fprintf(l.output, " %s", string(fieldsJSON))
			}
		}
		fmt.Fprintln(l.output)
	}
}

// DEBUG logs a debug message using MCode
func (l *Logger) DEBUG(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(DEBUG, mcode, optionalMessage, f)
}

// INFO logs an info message using MCode
func (l *Logger) INFO(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(INFO, mcode, optionalMessage, f)
}

// WARN logs a warning message using MCode
func (l *Logger) WARN(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WARN, mcode, optionalMessage, f)
}

// ERROR logs an error message using MCode
func (l *Logger) ERROR(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(ERROR, mcode, optionalMessage, f)
}

// FATAL logs a fatal message using MCode and exits
func (l *Logger) FATAL(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(FATAL, mcode, optionalMessage, f)
	os.Exit(1)
}
