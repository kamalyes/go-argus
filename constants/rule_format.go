/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\constants\rule_format.go
 * @Description: 格式校验规则名常量（网络、标识、编码等）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleEmail = "email" // 邮箱地址
	RuleE164  = "e164"  // E.164 电话号码

	RuleIP     = "ip"      // IP 地址（v4 或 v6）
	RuleIPAddr = "ip_addr" // IP 地址别名
	RuleIPv4   = "ipv4"   // IPv4 地址
	RuleIPv6   = "ipv6"   // IPv6 地址

	RuleCIDR   = "cidr"   // CIDR 表示法
	RuleCIDRv4 = "cidrv4" // IPv4 CIDR
	RuleCIDRv6 = "cidrv6" // IPv6 CIDR

	RuleMAC = "mac" // MAC 地址

	RuleHostname        = "hostname"         // 主机名
	RuleHostnameRFC1123 = "hostname_rfc1123" // RFC1123 主机名
	RuleFQDN            = "fqdn"             // 完全限定域名
	RuleHostnamePort    = "hostname_port"    // 主机名:端口
	RulePort            = "port"             // 端口号

	RuleURL        = "url"        // URL
	RuleURI        = "uri"        // URI
	RuleHTTPURL    = "http_url"   // HTTP URL
	RuleHTTPSURL   = "https_url"  // HTTPS URL
	RuleURLEncoded = "url_encoded" // URL 编码字符串

	RuleHTML        = "html"         // HTML 标签
	RuleHTMLEncoded = "html_encoded" // HTML 编码字符串

	RuleUUID         = "uuid"          // UUID（任意版本）
	RuleUUID3        = "uuid3"         // UUID v3
	RuleUUID4        = "uuid4"         // UUID v4
	RuleUUID5        = "uuid5"         // UUID v5
	RuleUUIDRFC4122  = "uuid_rfc4122"  // RFC4122 UUID
	RuleUUID3RFC4122 = "uuid3_rfc4122" // RFC4122 UUID v3
	RuleUUID4RFC4122 = "uuid4_rfc4122" // RFC4122 UUID v4
	RuleUUID5RFC4122 = "uuid5_rfc4122" // RFC4122 UUID v5

	RuleBase32       = "base32"       // Base32 编码
	RuleBase64       = "base64"       // Base64 编码
	RuleBase64URL    = "base64url"    // Base64URL 编码
	RuleBase64RawURL = "base64rawurl" // Base64RawURL 编码（无填充）

	RuleHexColor = "hexcolor" // 十六进制颜色值
	RuleRGB      = "rgb"      // RGB 颜色
	RuleRGBA     = "rgba"     // RGBA 颜色
	RuleHSL      = "hsl"      // HSL 颜色
	RuleHSLA     = "hsla"     // HSLA 颜色

	RuleSemver          = "semver"            // 语义化版本号
	RuleISBN10          = "isbn10"            // ISBN-10
	RuleISBN13          = "isbn13"            // ISBN-13
	RuleISSN            = "issn"              // ISSN
	RuleBIC             = "bic"               // BIC/SWIFT 代码
	RuleCron            = "cron"              // Cron 表达式
	RuleDataURI         = "datauri"           // Data URI
	RuleBCP47           = "bcp47"             // BCP47 语言标签
	RuleEthAddr         = "eth_addr"          // 以太坊地址
	RuleBtcAddr         = "btc_addr"          // 比特币地址

	RuleDatetime = "datetime" // 日期时间格式
	RuleTimezone = "timezone"  // 时区

	RuleFile     = "file"     // 文件路径
	RuleFilepath = "filepath" // 文件路径（别名）
	RuleDir      = "dir"      // 目录路径
	RuleDirpath  = "dirpath"  // 目录路径（别名）

	RuleMongoDB = "mongodb" // MongoDB ObjectID

	RuleLuhnChecksum    = "luhn_checksum"     // Luhn 校验和
	RuleCreditCard      = "credit_card"       // 信用卡号
	RuleDNSRFC1035Label = "dns_rfc1035_label" // DNS RFC1035 标签
)
