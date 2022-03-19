package constants

// resources
const (
    NODE_CERTS_DIR  = "data/certs/%s/"
    NODE_CERT_KEY   = NODE_CERTS_DIR + "client.key"
    NODE_CERT_CERT  = NODE_CERTS_DIR + "client.pem"

    DEFAULT_TEMPLATE_DIR = "scripts/templates/"
    DESIRE_TEMPLATE      = "sys_desire.json"
    REPORT_TEMPLATE      = "sys_report.json"
    APP_MYSQL_TEMPLATE   = "mysql.json"

    NODE_NAME = "nodeName"
    KEY_SRV_NAME = "core"
)

// app,template variables
const (
    BaetylInit       = "baetyl-init"
    BaetylCore       = "baetyl-core"
    BaetylBroker     = "baetyl-broker"
    BaetylFunction   = "baetyl-function"
    BaetylRule       = "baetyl-rule"
    BaetylAgent      = "baetyl-agent"
    BaetylLog        = "baetyl-log"
    BaetylGPUMetrics = "baetyl-gpu-metrics"

    TmplVarBaetylInit       = "BaetylInit"
    TmplVarBaetylCore       = "BaetylCore"
    TmplVarBaetylBroker     = "BaetylBroker"

    TmplVarBaetylInitVer       = "BaetylInitVersion"
    TmplVarBaetylCoreVer       = "BaetylCoreVersion"
    TmplVarBaetylBrokerVer     = "BaetylBrokerVersion"
)
