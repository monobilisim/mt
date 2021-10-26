

* [Yapılandırma dosyası](/config/config_sample.yml) /etc/mt/config.yml konumuna kaydedilir.
* Uygulamanın son sürümü, https://gitlab.com/monobilisim/mt/-/releases sayfasından indirilip /usr/local/bin/mt konumuna kaydedilir ve çalıştırılabilir yapılır.
* Örnek olarak, Asterisk sunucuları için görüşme kayıtlarını yedeklemek üzere aşağıdaki komut cron görevi olarak eklenir:
    ```
    /usr/local/bin/mt upload --source=/var/spool/asterisk/monitor/$(date -d yesterday +\%Y/\%m/\%d) --destination=minio/monitor/$(date -d yesterday +\%Y/\%m)/ -r -l debug -md5 -n -soe -rm
    ```
  * Genel parametreler:
    * `--config` / `-c`: Yapılandırma dosyası konumu. (Kullanılmazsa /etc/mt/config.yml konumuna bakılır.)
    * `--log-level` / `-l`: Log seviyesi (debug, info, warn veya error). (Kullanılmazsa yapılandırma dosyasındaki log seviyesine göre loglama yapılır.)
  * `upload` alt komutu için kullanılan parametrelerin açıklamaları:
    * `--source` / `-s`: Kaynak dosya veya dizin.
    * `--destination` / `-d`: Hedef.
      * Format: `<Yapılandırma dosyasındaki sunucu girdisi>/<bucket>/<prefix>`
    * `--recursive` / `-r`: Dizinleri taşımak için gerekli.
    * `--md5-validation` / `-md5`: Bu parametre kullanıldığında yüklenen dosya için MD5 kontrolü yapılır.
    * `--remove-source-files` / `-rm`: Bu parametre kullanıldığında yükleme ve kontroller başarılı ise kaynak dosya silinir.
    * `--stop-on-error` / `-soe`: Bu parametre kullanıldığında ilk hatada işlem sonlandırılır.
    * `--notify-errors` / `-n`: Bu parametre kullanıldığında hatalar için Rocket.Chat'e bildirim yapılır.
