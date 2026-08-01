/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-01 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-01 00:00:00
 * @FilePath: \go-argus\validate\image.go
 * @Description: 图片格式校验能力，基于魔数嗅探与尾部完整性检查支持 PNG、JPEG、GIF、WebP、BMP、TIFF、ICO、AVIF、HEIC、SVG
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"encoding/binary"
	"math/bits"

	"github.com/kamalyes/go-argus/i18n"
)

// 支持的图片格式规范名
const (
	imageFmtPNG  = "png"
	imageFmtJPEG = "jpeg"
	imageFmtGIF  = "gif"
	imageFmtWebP = "webp"
	imageFmtBMP  = "bmp"
	imageFmtTIFF = "tiff"
	imageFmtICO  = "ico"
	imageFmtAVIF = "avif"
	imageFmtHEIC = "heic"
	imageFmtSVG  = "svg"
)

// 图片格式位掩码，位序与 imageFormatNames 的索引一一对应，嗅探分发与白名单判定统一走位运算
const (
	bitPNG  = 1 << 0
	bitJPEG = 1 << 1
	bitGIF  = 1 << 2
	bitWebP = 1 << 3
	bitBMP  = 1 << 4
	bitTIFF = 1 << 5
	bitICO  = 1 << 6
	bitAVIF = 1 << 7
	bitHEIC = 1 << 8
	bitSVG  = 1 << 9
)

// imageFormatNames 格式位到规范格式名的映射表，索引即格式位中置位比特的位置
var imageFormatNames = [...]string{
	imageFmtPNG, imageFmtJPEG, imageFmtGIF, imageFmtWebP, imageFmtBMP,
	imageFmtTIFF, imageFmtICO, imageFmtAVIF, imageFmtHEIC, imageFmtSVG,
}

// imageFormatName 返回单个已知格式位对应的规范格式名，非单个已知位时返回空字符串
func imageFormatName(bit uint16) string {
	if bit == 0 || bit&(bit-1) != 0 {
		return ""
	}
	if n := bits.TrailingZeros16(bit); n < len(imageFormatNames) {
		return imageFormatNames[n]
	}
	return ""
}

// 魔数与尾部标记的机器字表示，读取后单次比较替代逐字节循环
const (
	pngMagicWord = 0x89504e470d0a1a0a // "\x89PNG\r\n\x1a\n" 大端
	riffWord     = 0x46464952         // "RIFF" 小端
	webpWord     = 0x50424557         // "WEBP" 小端
	tiffLEWord   = 0x002a4949         // "II\x2a\x00" 小端
	tiffBEWord   = 0x2a004d4d         // "MM\x00\x2a" 小端
	icoWord      = 0x00010000         // "\x00\x00\x01\x00" 小端
	ftypWord     = 0x66747970         // "ftyp" 大端
	brandAVIF    = 0x61766966         // "avif"
	brandAVIS    = 0x61766973         // "avis"
	brandHEIC    = 0x68656963         // "heic"
	brandHEIX    = 0x68656978         // "heix"
	brandHEVC    = 0x68657663         // "hevc"
	brandHEVX    = 0x68657678         // "hevx"
	brandHEIF    = 0x68656966         // "heif"
	brandMIF1    = 0x6d696631         // "mif1"
	brandMSF1    = 0x6d736631         // "msf1"
	iendWord     = 0x49454e44         // "IEND"
	iendCRCWord  = 0xae426082         // IEND CRC

	fold3 = 0x20202000 // 三字节格式名 token 的大小写归一掩码（低位填充字节不归一）
	fold4 = 0x20202020 // 四字节格式名 token 的大小写归一掩码
)

// svgScanLimit 限制 SVG 文本嗅探只扫描文件头部，避免大文本的全文扫描开销
const svgScanLimit = 512

// DetectImageFormat 基于魔数嗅探图片格式，返回规范格式名，无法识别时返回空字符串
// 首字节分派将顺序魔数链收敛为单次跳转，仅检查魔数不校验完整性，完整性检查由各 Is 系列函数或 imageIntact 执行
func DetectImageFormat(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	switch data[0] {
	case 0x89:
		if isPNGMagic(data) {
			return imageFmtPNG
		}
	case 0xff:
		if isJPEGMagic(data) {
			return imageFmtJPEG
		}
	case 'G':
		if isGIFMagic(data) {
			return imageFmtGIF
		}
	case 'R':
		if isWebPMagic(data) {
			return imageFmtWebP
		}
	case 'B':
		if isBMPMagic(data) {
			return imageFmtBMP
		}
	case 'I', 'M':
		if isTIFFMagic(data) {
			return imageFmtTIFF
		}
	case 0x00:
		if isICOMagic(data) {
			return imageFmtICO
		}
	}
	if bit := boxBrandBit(data); bit != 0 {
		return imageFormatName(bit)
	}
	if isSVGLeadByte(data[0]) && IsSVG(data) {
		return imageFmtSVG
	}
	return ""
}

// isSVGLeadByte 判断首字节是否可能开启 SVG 文档（BOM、空白或标签起始），用于跳过无谓的 SVG 扫描
func isSVGLeadByte(c byte) bool {
	switch c {
	case ' ', '\t', '\r', '\n', '<', 0xef:
		return true
	}
	return false
}

// NormalizeImageFormat 归一化图片格式名，忽略大小写且 jpg 等价 jpeg，未知格式名返回 false
// 与白名单解析共用 imageFormatTokenBit 的机器字归一比较，避免两套匹配逻辑并存
func NormalizeImageFormat(name string) (string, bool) {
	bit, ok := imageFormatTokenBit(name)
	if !ok {
		return "", false
	}
	return imageFormatName(bit), true
}

// IsImage 判断字节内容是否为有效且完整的图片（任意支持格式）
func IsImage(data []byte) bool {
	format := DetectImageFormat(data)
	return format != "" && imageIntact(format, data)
}

// IsImageInFormats 判断字节内容是否为白名单内的有效图片
// formats 为空格分隔的格式名，出现未知格式名时直接失败（fail-closed）
func IsImageInFormats(data []byte, formats string) bool {
	mask, ok := imageAllowedMask(formats)
	if !ok || mask == 0 {
		return false
	}
	format := DetectImageFormat(data)
	if format == "" {
		return false
	}
	bit, _ := imageFormatTokenBit(format)
	return mask&bit != 0 && imageIntact(format, data)
}

// IsPNG 判断字节内容是否为完整的 PNG 图片，魔数后必须存在 IEND 尾块
func IsPNG(data []byte) bool {
	return isPNGMagic(data) && pngTailOK(data)
}

// IsJPEG 判断字节内容是否为完整的 JPEG 图片，后缀必须为 EOI 标记以排除中缀含 EOI 的截断文件
func IsJPEG(data []byte) bool {
	return isJPEGMagic(data) && jpegTailOK(data)
}

// IsGIF 判断字节内容是否为完整的 GIF 图片，尾部必须为流结束符 0x3B
func IsGIF(data []byte) bool {
	return isGIFMagic(data) && gifTailOK(data)
}

// IsWebP 判断字节内容是否为有效的 WebP 图片，RIFF 容器内必须为 VP8/VP8L/VP8X 载荷
func IsWebP(data []byte) bool {
	return isWebPMagic(data)
}

// IsBMP 判断字节内容是否为有效的 BMP 图片，DIB 头大小必须为已知合法值
func IsBMP(data []byte) bool {
	return isBMPMagic(data)
}

// IsTIFF 判断字节内容是否为有效的 TIFF 图片，支持 II 与 MM 两种字节序
func IsTIFF(data []byte) bool {
	return isTIFFMagic(data)
}

// IsICO 判断字节内容是否为有效的 ICO 图标，目录项数量必须大于 0
func IsICO(data []byte) bool {
	return isICOMagic(data)
}

// IsAVIF 判断字节内容是否为有效的 AVIF 图片，基于 ftyp 容器的 major brand 判定
func IsAVIF(data []byte) bool {
	return boxBrandBit(data) == bitAVIF
}

// IsHEIC 判断字节内容是否为有效的 HEIC/HEIF 图片，基于 ftyp 容器的 major brand 判定
func IsHEIC(data []byte) bool {
	return boxBrandBit(data) == bitHEIC
}

// IsSVG 判断字节内容是否为有效的 SVG 文档
// 状态机嗅探：跳过 BOM、空白、<?xml…?>、<!--…-->、<!DOCTYPE…?> 后必须紧跟 <svg，全部扫描限制在文件头部
func IsSVG(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	// UTF-16 BOM 明确不支持
	if (data[0] == 0xff && data[1] == 0xfe) || (data[0] == 0xfe && data[1] == 0xff) {
		return false
	}
	limit := len(data)
	if limit > svgScanLimit {
		limit = svgScanLimit
	}
	i := 0
	// UTF-8 BOM
	if data[0] == 0xef && data[1] == 0xbb && data[2] == 0xbf {
		i = 3
	}
	for i < limit {
		switch data[i] {
		case ' ', '\t', '\r', '\n':
			i++
		case '<':
			switch {
			case matchAt(data, i, limit, "<svg"):
				return true
			case matchAt(data, i, limit, "<?xml"):
				end := indexByteWithin(data, i, limit, '>')
				if end < 0 {
					return false
				}
				i = end + 1
			case matchAt(data, i, limit, "<!--"):
				end := indexStrWithin(data, i+4, limit, "-->")
				if end < 0 {
					return false
				}
				i = end + 3
			case matchAt(data, i, limit, "<!DOCTYPE"):
				end := indexByteWithin(data, i, limit, '>')
				if end < 0 {
					return false
				}
				i = end + 1
			default:
				return false
			}
		default:
			return false
		}
	}
	return false
}

// ValidateImage 校验字节内容是否为有效且完整的图片
func ValidateImage(data []byte) CompareResult {
	format := DetectImageFormat(data)
	result := CompareResult{Actual: format, Expect: "valid image"}
	if format == "" || !imageIntact(format, data) {
		result.Message = i18n.Msg(MsgFormatImageInvalid)
		return result
	}
	result.Success = true
	return result
}

// ValidateImageInFormats 校验字节内容是否为白名单内的有效图片
// 白名单非法时快速失败，不再对字节内容执行嗅探
func ValidateImageInFormats(data []byte, formats string) CompareResult {
	result := CompareResult{Expect: formats}
	mask, ok := imageAllowedMask(formats)
	if !ok || mask == 0 {
		result.Message = i18n.Msg(MsgFormatImageNotAllowed, map[string]string{"actual": "", "expected": formats})
		return result
	}
	format := DetectImageFormat(data)
	result.Actual = format
	bit, _ := imageFormatTokenBit(format)
	if bit == 0 || mask&bit == 0 || !imageIntact(format, data) {
		result.Message = i18n.Msg(MsgFormatImageNotAllowed, map[string]string{"actual": format, "expected": formats})
		return result
	}
	result.Success = true
	return result
}

// imageIntact 按嗅探出的格式执行剩余完整性检查（魔数已由 DetectImageFormat 确认，不再重复嗅探）
func imageIntact(format string, data []byte) bool {
	switch format {
	case imageFmtPNG:
		return pngTailOK(data)
	case imageFmtJPEG:
		return jpegTailOK(data)
	case imageFmtGIF:
		return gifTailOK(data)
	case imageFmtWebP, imageFmtBMP, imageFmtTIFF, imageFmtICO,
		imageFmtAVIF, imageFmtHEIC, imageFmtSVG:
		return true
	}
	return false
}

// imageAllowedMask 将空格分隔的格式白名单解析为格式位掩码，出现未知格式名时返回 false
func imageAllowedMask(formats string) (uint16, bool) {
	mask := uint16(0)
	for i := 0; i < len(formats); {
		for i < len(formats) && (formats[i] == ' ' || formats[i] == '\t') {
			i++
		}
		start := i
		for i < len(formats) && formats[i] != ' ' && formats[i] != '\t' {
			i++
		}
		if i == start {
			break
		}
		bit, ok := imageFormatTokenBit(formats[start:i])
		if !ok {
			return 0, false
		}
		mask |= bit
	}
	return mask, true
}

// imageFormatTokenBit 单次机器字比较完成格式名归一化与位掩码映射
// token 按 |0x20 归一大小写（小写字母的原像只有大小写两种，无碰撞风险）
func imageFormatTokenBit(token string) (uint16, bool) {
	switch len(token) {
	case 3:
		word := uint32(token[0])<<24 | uint32(token[1])<<16 | uint32(token[2])<<8 | fold3
		switch word {
		case 0x706e6700: // png
			return bitPNG, true
		case 0x6a706700: // jpg
			return bitJPEG, true
		case 0x67696600: // gif
			return bitGIF, true
		case 0x626d7000: // bmp
			return bitBMP, true
		case 0x69636f00: // ico
			return bitICO, true
		case 0x73766700: // svg
			return bitSVG, true
		}
	case 4:
		word := uint32(token[0]) | uint32(token[1])<<8 | uint32(token[2])<<16 | uint32(token[3])<<24 | fold4
		switch word {
		case 0x6765706a: // jpeg
			return bitJPEG, true
		case 0x70626577: // webp
			return bitWebP, true
		case 0x66666974: // tiff
			return bitTIFF, true
		case 0x66697661: // avif
			return bitAVIF, true
		case 0x63696568: // heic
			return bitHEIC, true
		}
	}
	return 0, false
}

// isPNGMagic 单次 uint64 比较判定 PNG 魔数
func isPNGMagic(data []byte) bool {
	return len(data) >= 8 && binary.BigEndian.Uint64(data) == pngMagicWord
}

// isJPEGMagic 判定 JPEG SOI 魔数（\xff\xd8\xff）
func isJPEGMagic(data []byte) bool {
	return len(data) >= 3 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff
}

// isGIFMagic 判定 GIF87a/GIF89a 魔数
func isGIFMagic(data []byte) bool {
	if len(data) < 6 || data[0] != 'G' || data[1] != 'I' || data[2] != 'F' || data[3] != '8' {
		return false
	}
	return (data[4] == '7' || data[4] == '9') && data[5] == 'a'
}

// isWebPMagic 判定 RIFF 容器是否承载 WebP 载荷，VP 前缀覆盖 VP8/VP8L/VP8X
func isWebPMagic(data []byte) bool {
	return len(data) >= 14 &&
		binary.LittleEndian.Uint32(data) == riffWord &&
		binary.LittleEndian.Uint32(data[8:]) == webpWord &&
		data[12] == 'V' && data[13] == 'P'
}

// isBMPMagic 判定 BMP 魔数且 DIB 头大小合法
func isBMPMagic(data []byte) bool {
	if len(data) < 18 || data[0] != 'B' || data[1] != 'M' {
		return false
	}
	switch binary.LittleEndian.Uint32(data[14:18]) {
	case 12, 40, 52, 56, 64, 108, 124:
		return true
	}
	return false
}

// isTIFFMagic 单次 uint32 比较判定 TIFF 双字节序魔数
func isTIFFMagic(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	word := binary.LittleEndian.Uint32(data)
	return word == tiffLEWord || word == tiffBEWord
}

// isICOMagic 判定 ICO 魔数且图标目录项数量大于 0
func isICOMagic(data []byte) bool {
	return len(data) >= 6 &&
		binary.LittleEndian.Uint32(data) == icoWord &&
		(data[4] != 0 || data[5] != 0)
}

// boxBrandBit 嗅探 ftyp 容器的图片格式位，非 AVIF/HEIC 容器返回 0
func boxBrandBit(data []byte) uint16 {
	if len(data) < 12 || binary.BigEndian.Uint32(data[4:8]) != ftypWord {
		return 0
	}
	switch binary.BigEndian.Uint32(data[8:12]) {
	case brandAVIF, brandAVIS:
		return bitAVIF
	case brandHEIC, brandHEIX, brandHEVC, brandHEVX, brandHEIF, brandMIF1, brandMSF1:
		return bitHEIC
	}
	return 0
}

// pngTailOK 三次 uint32 比较判定 PNG IEND 尾块（4 字节零长度 + "IEND" + CRC）
func pngTailOK(data []byte) bool {
	n := len(data)
	if n < 12 {
		return false
	}
	return binary.BigEndian.Uint32(data[n-12:n-8]) == 0 &&
		binary.BigEndian.Uint32(data[n-8:n-4]) == iendWord &&
		binary.BigEndian.Uint32(data[n-4:]) == iendCRCWord
}

// jpegTailOK 判定 JPEG 尾部为 EOI 标记
func jpegTailOK(data []byte) bool {
	n := len(data)
	return n >= 4 && data[n-2] == 0xff && data[n-1] == 0xd9
}

// gifTailOK 判定 GIF 尾部为流结束符
func gifTailOK(data []byte) bool {
	return len(data) >= 7 && data[len(data)-1] == 0x3b
}

// matchAt 判断 data 从 i 开始（不越过 limit）是否匹配给定字符串
func matchAt(data []byte, i, limit int, s string) bool {
	if i+len(s) > limit {
		return false
	}
	return bytesEqualStr(data[i:i+len(s)], s)
}

// bytesEqualStr 零分配判断字节内容与字符串是否逐字节相等
func bytesEqualStr(data []byte, s string) bool {
	if len(data) != len(s) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if data[i] != s[i] {
			return false
		}
	}
	return true
}

// indexByteWithin 在 [from, limit) 范围内查找字节，未命中返回 -1
func indexByteWithin(data []byte, from, limit int, delim byte) int {
	for i := from; i < limit; i++ {
		if data[i] == delim {
			return i
		}
	}
	return -1
}

// indexStrWithin 在 [from, limit) 范围内查找字符串，未命中返回 -1
func indexStrWithin(data []byte, from, limit int, s string) int {
	for i := from; i+len(s) <= limit; i++ {
		if bytesEqualStr(data[i:i+len(s)], s) {
			return i
		}
	}
	return -1
}
