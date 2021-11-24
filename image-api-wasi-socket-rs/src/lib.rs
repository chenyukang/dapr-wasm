use image::io::Reader as ImageReader;
use std::io::Cursor;

pub fn detect_qr(bytes: &[u8]) -> Result<String, Box<dyn std::error::Error>> {
    let img = ImageReader::new(Cursor::new(bytes))
        .with_guessed_format()?
        .decode()?;

    // Use default decoder
    let decoder = bardecoder::default_decoder();

    let results = decoder.decode(&img);
    if results.is_empty() || results[0].is_err() {
        return Err("No QR code found".into());
    }
    Ok("Valid QR code".into())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_detect_qr() {
        let bytes = include_bytes!("../test/images/wechat.png");
        let result = detect_qr(bytes);
        assert!(result.is_err());
    }

    #[test]
    fn test_detect_qr_valid() {
        let bytes = include_bytes!("../test/images/qr.jpg");
        let result = detect_qr(bytes);
        assert!(result.is_ok());
    }
}
