/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-01 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-01 00:00:00
 * @FilePath: \go-argus\validate\image_test.go
 * @Description: 图片魔数嗅探与完整性校验单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"strings"
	"testing"
)

// 最小合法样本：只保证魔数与尾部结构真实，字段内容从简
const (
	pngSample = "\x89PNG\r\n\x1a\n" + // PNG 魔数
		"\x00\x00\x00\x0dIHDR" + // IHDR 块头（长度 13）
		"\x08\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00" + // IHDR 数据（宽1高1、RGBA）
		"\x00\x00\x00\x00" + // IHDR CRC（校验不解析内容）
		"\x00\x00\x00\x00IEND\xae\x42\x60\x82" // IEND 尾块
	jpegSample  = "\xff\xd8\xff\xe0\x00\x10JFIF\x00\x01\x01\x00\x00\x01\x00\x01\x00\x00\xff\xd9" // SOI + JFIF 段 + EOI
	jpegFakeEOI = "\xff\xd8\xff\xe0\x00\x10JFIF\x00\x01\x01\x00\x00\x01\x00\x01\x00\x00\xff\xd9TRUNC"
	gifSample   = "GIF89a\x01\x00\x01\x00\x00\x00\x00\x00\x3b" // 逻辑屏幕描述符 + 流结束符
	gif87Sample = "GIF87a\x01\x00\x01\x00\x00\x00\x00\x00\x3b"
	webpSample  = "RIFF\x24\x00\x00\x00WEBPVP8L\x00\x00\x00\x2f\x00\x00\x00"
	webpEmpty   = "RIFF\x24\x00\x00\x00WEBPADMP" // RIFF 容器但非 VP 载荷
	bmpSample   = "BM\x42\x00\x00\x00\x00\x00\x00\x00\x36\x00\x00\x00\x28\x00\x00\x00\x00\x00\x00\x00"
	bmpBadDIB   = "BM\x42\x00\x00\x00\x00\x00\x00\x00\x36\x00\x00\x00\x30\x00\x00\x00\x00\x00\x00\x00" // DIB 头大小非法
	icoSample   = "\x00\x00\x01\x00\x01\x00\x30\x30\x10\x10\x00\x00\x01\x00\x20\x00"
	icoEmpty    = "\x00\x00\x01\x00\x00\x00\x30\x30\x10\x10\x00\x00\x01\x00\x20\x00" // 图标数为 0
	tiffSample  = "II\x2a\x00\x08\x00\x00\x00\x00\x00"
	tiffSample2 = "MM\x00\x2a\x00\x00\x00\x08\x00\x00"
	avifSample  = "\x00\x00\x00\x20ftypavif\x00\x00\x00\x00mdat\x00\x00\x00\x00"
	avifAvis    = "\x00\x00\x00\x20ftypavis\x00\x00\x00\x00mdat\x00\x00\x00\x00"
	heicSample  = "\x00\x00\x00\x20ftypheic\x00\x00\x00\x00mdat\x00\x00\x00\x00"
	heicMif1    = "\x00\x00\x00\x20ftypmif1\x00\x00\x00\x00mdat\x00\x00\x00\x00"
)

// svg 样本集中覆盖前缀变体与误报场景
var (
	svgPlain   = []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`)
	svgXMLDecl = []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<svg/>`)
	svgBOM     = append([]byte("\xef\xbb\xbf"), []byte(`<svg/>`)...)
	svgComment = []byte(`<!-- banner -->` + "\n" + `<svg/>`)
	svgDoctype = []byte(`<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">` + "\n" + `<svg/>`)
	svgHTML    = []byte(`<!DOCTYPE html><html><body><svg/></body></html>`)
	svgUTF16   = []byte("\xff\xfe<\x00s\x00v\x00g\x00")
	svgText    = []byte(`hello world <svg/>`)
	svgLongXML = append([]byte(`<?xml version="1.0" `), append([]byte(strings.Repeat("a", 520)), []byte(`?><svg/>`)...)...)
)

func TestDetectImageFormat(t *testing.T) {
	cases := []struct {
		name   string
		data   []byte
		expect string
	}{
		{"png", []byte(pngSample), "png"},
		{"png_magic_only", append([]byte(pngSample), 'x'), "png"},
		{"jpeg", []byte(jpegSample), "jpeg"},
		{"jpeg_fake_eoi", []byte(jpegFakeEOI), "jpeg"},
		{"gif89", []byte(gifSample), "gif"},
		{"gif87", []byte(gif87Sample), "gif"},
		{"webp", []byte(webpSample), "webp"},
		{"webp_empty", []byte(webpEmpty), ""},
		{"bmp", []byte(bmpSample), "bmp"},
		{"bmp_bad_dib", []byte(bmpBadDIB), ""},
		{"tiff_le", []byte(tiffSample), "tiff"},
		{"tiff_be", []byte(tiffSample2), "tiff"},
		{"ico", []byte(icoSample), "ico"},
		{"ico_empty", []byte(icoEmpty), ""},
		{"avif", []byte(avifSample), "avif"},
		{"avif_avis", []byte(avifAvis), "avif"},
		{"heic", []byte(heicSample), "heic"},
		{"heic_mif1", []byte(heicMif1), "heic"},
		{"svg_plain", svgPlain, "svg"},
		{"svg_xml_decl", svgXMLDecl, "svg"},
		{"svg_bom", svgBOM, "svg"},
		{"svg_comment", svgComment, "svg"},
		{"svg_doctype", svgDoctype, "svg"},
		{"svg_html", svgHTML, ""},
		{"svg_utf16", svgUTF16, ""},
		{"svg_text_prefix", svgText, ""},
		{"svg_long_xml", svgLongXML, ""},
		{"empty", nil, ""},
		{"random", []byte("not an image at all"), ""},
	}
	for _, c := range cases {
		if got := DetectImageFormat(c.data); got != c.expect {
			t.Errorf("%s: DetectImageFormat = %q, want %q", c.name, got, c.expect)
		}
	}
}

func TestImageIntegrity(t *testing.T) {
	pngTruncated := []byte(pngSample[:len(pngSample)-6])
	gifTruncated := []byte("GIF89a\x01\x00\x01\x00\x00\x00\x00\x00")
	pngBadMagic := []byte("\x89PNX\r\n\x1a\n" + pngSample[8:])
	jpegBadMagic := append([]byte("\xff\xd8\xfe"), jpegSample[3:]...)

	cases := []struct {
		name   string
		data   []byte
		fn     func([]byte) bool
		expect bool
	}{
		{"png_ok", []byte(pngSample), IsPNG, true},
		{"png_truncated", pngTruncated, IsPNG, false},
		{"png_bad_magic", pngBadMagic, IsPNG, false},
		{"png_empty", nil, IsPNG, false},
		{"jpeg_ok", []byte(jpegSample), IsJPEG, true},
		{"jpeg_fake_eoi_suffix", []byte(jpegFakeEOI), IsJPEG, false},
		{"jpeg_truncated", []byte(jpegSample[:len(jpegSample)-3]), IsJPEG, false},
		{"jpeg_bad_magic", jpegBadMagic, IsJPEG, false},
		{"gif_ok", []byte(gifSample), IsGIF, true},
		{"gif87_ok", []byte(gif87Sample), IsGIF, true},
		{"gif_truncated", gifTruncated, IsGIF, false},
		{"gif_bad_magic", []byte("GIF88a\x01\x00\x01\x00\x00\x00\x00\x3b"), IsGIF, false},
		{"webp_ok", []byte(webpSample), IsWebP, true},
		{"webp_empty_payload", []byte(webpEmpty), IsWebP, false},
		{"webp_short", []byte("RIFF\x00\x00\x00\x00WEBPV"), IsWebP, false},
		{"bmp_ok", []byte(bmpSample), IsBMP, true},
		{"bmp_bad_dib", []byte(bmpBadDIB), IsBMP, false},
		{"bmp_short", []byte("BM\x28\x00"), IsBMP, false},
		{"tiff_ok", []byte(tiffSample), IsTIFF, true},
		{"tiff_be_ok", []byte(tiffSample2), IsTIFF, true},
		{"tiff_bad", []byte("II\x2b\x00"), IsTIFF, false},
		{"ico_ok", []byte(icoSample), IsICO, true},
		{"ico_empty", []byte(icoEmpty), IsICO, false},
		{"avif_ok", []byte(avifSample), IsAVIF, true},
		{"avif_avis", []byte(avifAvis), IsAVIF, true},
		{"avif_no", []byte(heicSample), IsAVIF, false},
		{"heic_ok", []byte(heicSample), IsHEIC, true},
		{"heic_mif1", []byte(heicMif1), IsHEIC, true},
		{"heic_no", []byte(avifSample), IsHEIC, false},
		{"svg_plain", svgPlain, IsSVG, true},
		{"svg_xml_decl", svgXMLDecl, IsSVG, true},
		{"svg_bom", svgBOM, IsSVG, true},
		{"svg_comment", svgComment, IsSVG, true},
		{"svg_doctype", svgDoctype, IsSVG, true},
		{"svg_html", svgHTML, IsSVG, false},
		{"svg_utf16", svgUTF16, IsSVG, false},
		{"svg_text_prefix", svgText, IsSVG, false},
		{"svg_long_xml", svgLongXML, IsSVG, false},
	}
	for _, c := range cases {
		if got := c.fn(c.data); got != c.expect {
			t.Errorf("%s: got %v, want %v", c.name, got, c.expect)
		}
	}
}

func TestIsImage(t *testing.T) {
	valid := [][]byte{
		[]byte(pngSample),
		[]byte(jpegSample),
		[]byte(gifSample),
		[]byte(webpSample),
		[]byte(bmpSample),
		[]byte(tiffSample),
		[]byte(icoSample),
		[]byte(avifSample),
		[]byte(heicSample),
		svgPlain,
		svgXMLDecl,
	}
	for i, data := range valid {
		if !IsImage(data) {
			t.Errorf("valid[%d]: IsImage(%q…) = false, want true", i, data[:min(8, len(data))])
		}
	}

	invalid := [][]byte{
		nil,
		{},
		[]byte(pngSample[:len(pngSample)-6]), // 截断 PNG
		[]byte(jpegFakeEOI),                  // 中缀 EOI 的截断 JPEG
		[]byte("GIF89a\x01\x00\x01\x00\x00\x00\x00"), // 截断 GIF
		[]byte(webpEmpty),         // 空 RIFF 载荷
		[]byte(bmpBadDIB),         // DIB 头非法
		[]byte(icoEmpty),          // 图标数为 0
		[]byte("plain text data"), // 纯文本
		svgHTML,                   // HTML 误报
	}
	for i, data := range invalid {
		if IsImage(data) {
			t.Errorf("invalid[%d]: IsImage = true, want false", i)
		}
	}
}

func TestIsImageInFormats(t *testing.T) {
	png := []byte(pngSample)
	jpeg := []byte(jpegSample)
	gif := []byte(gifSample)

	cases := []struct {
		name    string
		data    []byte
		formats string
		expect  bool
	}{
		{"png_in_png", png, "png", true},
		{"png_in_mixed", png, "jpg png", true},
		{"png_in_upper", png, "PNG", true},
		{"png_in_mixed_case", png, "Png WebP", true},
		{"png_not_in_jpeg", png, "jpeg", false},
		{"jpeg_alias_jpg", jpeg, "jpg", true},
		{"jpeg_alias_jpeg", jpeg, "jpeg", true},
		{"gif_multi_tabs", gif, "png\twebp\tgif", true},
		{"unknown_format_fail_closed", png, "png tiff2", false},
		{"empty_formats", png, "", false},
		{"spaces_only", png, "   ", false},
		{"truncated_png", png[:len(png)-6], "png", false},
	}
	for _, c := range cases {
		if got := IsImageInFormats(c.data, c.formats); got != c.expect {
			t.Errorf("%s: IsImageInFormats(%q) = %v, want %v", c.name, c.formats, got, c.expect)
		}
	}
}

func TestNormalizeImageFormat(t *testing.T) {
	valid := map[string]string{
		"png": "png", "PNG": "png", "Png": "png",
		"jpg": "jpeg", "JPG": "jpeg", "jpeg": "jpeg", "JPEG": "jpeg",
		"gif": "gif", "webp": "webp", "bmp": "bmp", "tiff": "tiff",
		"ico": "ico", "avif": "avif", "heic": "heic", "svg": "svg",
	}
	for in, want := range valid {
		got, ok := NormalizeImageFormat(in)
		if !ok || got != want {
			t.Errorf("NormalizeImageFormat(%q) = (%q, %v), want (%q, true)", in, got, ok, want)
		}
	}

	for _, in := range []string{"", "tiff2", "jp2", "jpgg", "raw"} {
		if got, ok := NormalizeImageFormat(in); ok {
			t.Errorf("NormalizeImageFormat(%q) = (%q, true), want false", in, got)
		}
	}
}

func TestValidateImageCompareResult(t *testing.T) {
	result := ValidateImage([]byte(pngSample))
	if !result.Success || result.Actual != "png" || result.Expect != "valid image" {
		t.Errorf("ValidateImage(png) = %+v, want success with actual png", result)
	}
	if result.Message != "" {
		t.Errorf("ValidateImage(png).Message = %q, want empty", result.Message)
	}

	result = ValidateImage([]byte("not an image"))
	if result.Success || result.Message == "" {
		t.Errorf("ValidateImage(text) = %+v, want failure with message", result)
	}

	result = ValidateImage(nil)
	if result.Success || result.Message == "" {
		t.Errorf("ValidateImage(nil) = %+v, want failure with message", result)
	}

	result = ValidateImage([]byte(pngSample[:len(pngSample)-6]))
	if result.Success || result.Message == "" {
		t.Errorf("ValidateImage(truncated png) = %+v, want failure with message", result)
	}
}

func TestValidateImageInFormatsCompareResult(t *testing.T) {
	result := ValidateImageInFormats([]byte(gifSample), "png jpg")
	if result.Success || result.Actual != "gif" || result.Expect != "png jpg" || result.Message == "" {
		t.Errorf("ValidateImageInFormats(gif, png jpg) = %+v, want failure listing actual gif", result)
	}

	result = ValidateImageInFormats([]byte(jpegSample), "webp jpg")
	if !result.Success || result.Actual != "jpeg" {
		t.Errorf("ValidateImageInFormats(jpeg, webp jpg) = %+v, want success with actual jpeg", result)
	}

	result = ValidateImageInFormats([]byte(pngSample), "png tiffx")
	if result.Success || result.Message == "" {
		t.Errorf("ValidateImageInFormats(png, png tiffx) = %+v, want failure for unknown format", result)
	}

	result = ValidateImageInFormats([]byte(pngSample), "")
	if result.Success || result.Message == "" {
		t.Errorf("ValidateImageInFormats(png, empty) = %+v, want failure for empty whitelist", result)
	}
}

func TestImageFormatBitConsistency(t *testing.T) {
	// 每个规范格式名的位掩码必须唯一且非零，且位掩码能还原回规范名
	formats := []string{"png", "jpeg", "gif", "webp", "bmp", "tiff", "ico", "avif", "heic", "svg"}
	seen := make(map[uint16]string, len(formats))
	for _, f := range formats {
		bit, ok := imageFormatTokenBit(f)
		if !ok || bit == 0 {
			t.Errorf("imageFormatTokenBit(%q) = (%d, %v), want non-zero ok", f, bit, ok)
			continue
		}
		if prev, dup := seen[bit]; dup {
			t.Errorf("imageFormatTokenBit(%q) duplicates %q", f, prev)
		}
		seen[bit] = f
		if name := imageFormatName(bit); name != f {
			t.Errorf("imageFormatName(%#x) = %q, want %q", bit, name, f)
		}
	}
	if bit, ok := imageFormatTokenBit("nope"); ok || bit != 0 {
		t.Errorf("imageFormatTokenBit(unknown) = (%d, %v), want (0, false)", bit, ok)
	}
	if name := imageFormatName(0); name != "" {
		t.Errorf("imageFormatName(0) = %q, want empty", name)
	}
	if name := imageFormatName(bitPNG | bitJPEG); name != "" {
		t.Errorf("imageFormatName(multi-bit) = %q, want empty", name)
	}
	if name := imageFormatName(1 << 15); name != "" {
		t.Errorf("imageFormatName(out-of-range bit) = %q, want empty", name)
	}
}

func BenchmarkDetectImageFormat(b *testing.B) {
	datasets := []struct {
		name string
		data []byte
	}{
		{"png", []byte(pngSample)},
		{"jpeg", []byte(jpegSample)},
		{"gif", []byte(gifSample)},
		{"webp", []byte(webpSample)},
		{"bmp", []byte(bmpSample)},
		{"tiff", []byte(tiffSample)},
		{"ico", []byte(icoSample)},
		{"avif", []byte(avifSample)},
		{"heic", []byte(heicSample)},
		{"svg", svgPlain},
		{"unknown", []byte("plain text data")},
	}
	for _, d := range datasets {
		b.Run(d.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = DetectImageFormat(d.data)
			}
		})
	}
}

func BenchmarkIsImage(b *testing.B) {
	data := []byte(pngSample)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsImage(data)
	}
}

func BenchmarkIsImageInFormats(b *testing.B) {
	data := []byte(pngSample)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsImageInFormats(data, "png jpg webp")
	}
}

func BenchmarkIsImageInFormatsParallel(b *testing.B) {
	data := []byte(pngSample)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = IsImageInFormats(data, "png jpg webp")
		}
	})
}

func BenchmarkNormalizeImageFormat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = NormalizeImageFormat("webp")
	}
}
