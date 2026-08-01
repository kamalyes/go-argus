/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-01 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-01 00:00:00
 * @FilePath: \go-argus\constants\rule_image.go
 * @Description: 图片格式校验规则名常量（PNG、JPEG、GIF、WebP、BMP、TIFF、ICO、AVIF、HEIC、SVG）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleImage = "image" // 图片（任意支持格式，可用空格分隔参数指定白名单）
	RulePNG   = "png"   // PNG 图片（含 IEND 尾块完整性检查）
	RuleJPG   = "jpg"   // JPEG 图片（RuleJPEG 别名）
	RuleJPEG  = "jpeg"  // JPEG 图片（含 EOI 尾部完整性检查）
	RuleGIF   = "gif"   // GIF 图片（含 0x3B 尾部完整性检查）
	RuleWebP  = "webp"  // WebP 图片
	RuleBMP   = "bmp"   // BMP 图片
	RuleTIFF  = "tiff"  // TIFF 图片
	RuleICO   = "ico"   // ICO 图标
	RuleAVIF  = "avif"  // AVIF 图片
	RuleHEIC  = "heic"  // HEIC/HEIF 图片
	RuleSVG   = "svg"   // SVG 文档
)
