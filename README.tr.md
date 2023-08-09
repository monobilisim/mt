[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![GPL License][license-shield]][license-url]

[![Readme in English](https://img.shields.io/badge/Readme-English-blue)](README.md)
[![Readme in Turkish](https://img.shields.io/badge/Readme-Turkish-red)](README.tr.md)

<div align="center"> 
<a href="https://monobilisim.com.tr/">
  <img src="https://monobilisim.com.tr/images/mono-bilisim.svg" width="340"/>
</a>

<h2 align="center">mt</h2>
<b>mt</b>, Asterisk sunucularının arama kayıtlarını yedeklemek ve bunları Minio dahil çeşitli hedeflere aktarmak için bir araçtır.
</div>

---

## İçindekiler

- [İçindekiler](#i̇çindekiler)
- [Kurulum](#kurulum)
- [Kullanım](#kullanım)
- [Parametreler](#parametreler)
  - [Genel Parametreler](#genel-parametreler)
  - [Upload Subcommand Parameters](#upload-subcommand-parameters)
- [License](#license)

---

## Kurulum

1. Yapılandırma dosyasını `/etc/mt/config.yml` olarak kaydedin veya [config_sample.yml](/config/config_sample.yml) dosyasını bir şablon olarak kullanın.

2. En son `mt` sürümünü [GitLab Sürümlerinden](https://github.com/monobilisim/mt/releases) indirin ve `/usr/local/bin/mt` dizinine yerleştirin. Çalıştırılabilir yapın.

---

## Kullanım

Asterisk sunucularının arama kayıtlarını yedeklemek ve Minio'ya aktarmak için:

```bash
/usr/local/bin/mt upload --kaynak=/var/spool/asterisk/monitor/$(date -d yesterday +\%Y/\%m/\%d) --hedef=minio/monitor/$(date -d yesterday +\%Y/\%m)/ -r -l debug -md5 -n -soe -rm -dmp
```

## Parametreler

### Genel Parametreler

- `--config` / `-c`: Yapılandırma dosyasının konumu. (Kullanılmazsa, varsayılan `/etc/mt/config.yml` olarak belirlenir.)
- `--log-level` / `-l`: Günlük seviyesi (debug, info, warn veya error). (Kullanılmazsa, günlük seviyesi yapılandırma dosyasına göre belirlenir.)

### Upload Subcommand Parameters

- `--source` / `-s`: Kaynak dosya veya dizini.
- `--destination` / `-d`: Hedef dosya veya dizin.
    - Format: `<Yapılandırmadaki Sunucu/<bucket>/<prefix>`
- `--recursive` / `-r`: Dizinleri taşımak için gerekli.
- `--md5-validation` / `-md5`: Kullanıldığında, yüklenen dosya için MD5 doğrulaması yapılır.
- `--remove-source-files` / `-rm`: Kullanıldığında, başarılı yükleme ve doğrulama sonrasında kaynak dosya silinir.
- `--stop-on-error` / `-soe`: Kullanıldığında, ilk hatada işlem sonlandırılır.
- `--notify-errors` / `-n`: Kullanıldığında, hatalar Rocket.Chat'e bildirilir.
- `--disable-multipart` / `-dmp`: Kullanıldığında, büyük dosyalar için çok parçalı yükleme devre dışı bırakılır.

---

## License

mt, GPL-3.0 lisanslıdır. Detaylar için [LICENSE](LICENSE) dosyasına bakınız.

[contributors-shield]: https://img.shields.io/github/contributors/monobilisim/mt.svg?style=for-the-badge
[contributors-url]: https://github.com/monobilisim/mt/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/monobilisim/mt.svg?style=for-the-badge
[forks-url]: https://github.com/monobilisim/mt/network/members
[stars-shield]: https://img.shields.io/github/stars/monobilisim/mt.svg?style=for-the-badge
[stars-url]: https://github.com/monobilisim/mt/stargazers
[issues-shield]: https://img.shields.io/github/issues/monobilisim/mt.svg?style=for-the-badge
[issues-url]: https://github.com/monobilisim/mt/issues
[license-shield]: https://img.shields.io/github/license/monobilisim/mt.svg?style=for-the-badge
[license-url]: https://github.com/monobilisim/mt/blob/master/LICENSE.txt