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
<b>mt</b> is a tool for backing up call recordings for Asterisk servers and transferring them to various destinations, including Minio.
</div>

---

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Installation](#installation)
- [Usage](#usage)
- [Parameters](#parameters)
  - [General Parameters](#general-parameters)
  - [Upload Subcommand Parameters](#upload-subcommand-parameters)
- [License](#license)

---

## Installation

1. Save the configuration file as `/etc/mt/config.yml` or use [config_sample.yml](/config/config_sample.yml) as a template.

2. Download the latest version of `mt` from [GitHub Releases](https://github.com/monobilisim/mt/releases) and place it in `/usr/local/bin/mt`. Make it executable.

---

## Usage

To backup call recordings for Asterisk servers and transfer them to Minio:

```bash
/usr/local/bin/mt upload --source=/var/spool/asterisk/monitor/$(date -d yesterday +\%Y/\%m/\%d) --destination=minio/monitor/$(date -d yesterday +\%Y/\%m)/ -r -l debug -md5 -n -soe -rm -dmp
```
## Parameters

### General Parameters

- `--config` / `-c`: Location of the configuration file. (If not used, the default is `/etc/mt/config.yml`.)
- `--log-level` / `-l`: Log level (debug, info, warn, or error). (If not used, log level is determined by the configuration file.)

### Upload Subcommand Parameters

- `--source` / `-s`: Source file or directory.
- `--destination` / `-d`: Destination.
    - Format: `<Server entry in the configuration>/<bucket>/<prefix>`
- `--recursive` / `-r`: Required to move directories.
- `--md5-validation` / `-md5`: When used, MD5 validation is performed for the uploaded file.
- `--remove-source-files` / `-rm`: When used, source file is deleted after successful upload and validation.
- `--stop-on-error` / `-soe`: When used, the process is terminated on the first error.
- `--notify-errors` / `-n`: When used, errors are notified to Rocket.Chat.
- `--disable-multipart` / `-dmp`: When used, multipart upload is disabled for large files.

---

## License

mt is GPL-3.0 licensed. See [LICENSE](LICENSE) file for details.


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