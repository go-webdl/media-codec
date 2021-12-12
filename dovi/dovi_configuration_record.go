package dovi

import (
	"encoding/binary"
	"io"
)

// 2.2 Dolby Vision configuration boxes
//
// https://professional.dolby.com/siteassets/content-creation/dolby-vision-for-content-creators/dolby_vision_bitstreams_within_the_iso_base_media_file_format_dec2017.pdf
//
// An ISO base media file that carries a Dolby Vision stream must contain boxes
// that signal the configuration information about the contained Dolby Vision
// stream. The configuration information is required to initialize the Dolby
// Vision decoder.
type DOVIDecoderConfigurationRecord struct {
	VersionMajor            uint8
	VersionMinor            uint8
	Profile                 uint8
	Level                   uint8
	RPUPresent              bool
	ELPresent               bool
	BLPresent               bool
	BLSignalCompatibilityID uint8
}

func (b *DOVIDecoderConfigurationRecord) RecordSize() (size uint32) {
	// unsigned int (8) dv_version_major;
	// unsigned int (8) dv_version_minor;
	// unsigned int (7) dv_profile;
	// unsigned int (6) dv_level;
	// bit (1) rpu_present_flag;
	// bit (1) el_present_flag;
	// bit (1) bl_present_flag;
	// unsigned int (4) dv_bl_signal_compatibility_id;
	// const unsigned int (28) reserved = 0;
	// const unsigned int (32)[4] reserved = 0;
	size = 24
	return
}

func (b *DOVIDecoderConfigurationRecord) RecordRead(r io.Reader) (err error) {
	var tmp [24]uint8
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return
	}
	b.VersionMajor = tmp[0]
	b.VersionMinor = tmp[1]
	b.Profile = tmp[2] >> 1
	b.Level = ((tmp[2] & 0b1) << 5) | ((tmp[3] & 0b11111000) >> 3)
	b.RPUPresent = (tmp[3] & 0b00000100) > 0
	b.ELPresent = (tmp[3] & 0b00000010) > 0
	b.BLPresent = (tmp[3] & 0b00000001) > 0
	b.BLSignalCompatibilityID = tmp[4] >> 4
	return
}

func (b *DOVIDecoderConfigurationRecord) RecordWrite(w io.Writer) (err error) {
	var tmp [24]uint8
	tmp[0] = b.VersionMajor
	tmp[1] = b.VersionMinor
	tmp[2] = (b.Profile << 1) | ((b.Level >> 5) & 0b1)
	tmp[3] = ((b.Level << 3) & 0b11111000)
	if b.RPUPresent {
		tmp[3] |= 0b00000100
	}
	if b.ELPresent {
		tmp[3] |= 0b00000010
	}
	if b.BLPresent {
		tmp[3] |= 0b00000001
	}
	tmp[4] = b.BLSignalCompatibilityID << 4
	if err = binary.Write(w, binary.BigEndian, &tmp); err != nil {
		return
	}
	return
}
