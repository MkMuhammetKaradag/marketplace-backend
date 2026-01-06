# Marketplace Backend

Bu proje, Go dili kullanılarak geliştirilmiştir.


### Servisler

- **API Gateway**: İstemcilerden gelen tüm isteklerin tek giriş noktasıdır. Yönlendirme, kimlik doğrulama (auth), rate limiting ve istek doğrulama işlemlerini gerçekleştirir. Fiber framework'ü kullanır.
- **Product Service**: Ürünlerin eklenmesi, güncellenmesi, listelenmesi ve aranması işlemlerinden sorumludur. Vektör tabanlı arama özellikleri için PostgreSQL (pgvector) kullanır.
- **User Service**: Kullanıcı kaydı, giriş işlemler, profil yönetimi ve kullanıcıyla ilgili diğer işlemleri yürütür.
- **Seller Service**: Satıcıların mağaza yönetimi, ürün envanteri ve satıcıya özgü operasyonları yönetir.

### İletişim

- **İç İletişim**: Servisler arası iletişim ağırlıklı olarak gRPC ve asenkron işlemler için Apache Kafka üzerinden sağlanır.
- **Dış İletişim**: İstemciler (Frontend, Mobile vb.) API Gateway ile RESTful HTTP üzerinden haberleşir.

## 🛠 Teknoloji Yığını

- **Programlama Dili**: Go (Golang)
- **Veritabanı**: PostgreSQL (pgvector eklentisi ile birlikte)
- **Önbellekleme (Cache)**: Redis
- **Mesaj Kuyruğu**: Apache Kafka (Zookeeper ile)
- **AI & Vektör Arama**: Ollama (`nomic-embed-text` modeli ile embedding işlemleri)
- **API Framework**: Fiber (Go için hızlı bir web framework)
- **RPC Framework**: gRPC
- **Konteynerleştirme**: Docker & Docker Compose

## 🚀 Kurulum ve Çalıştırma

Projenin yerel ortamda çalıştırılması için aşağıdaki adımları izleyebilirsiniz.

### Gereksinimler

- Go 1.21 veya üzeri
- Docker ve Docker Compose

### Altyapıyı Hazırlama

Veritabanı, Redis, Kafka ve Ollama gibi bağımlı servisleri Docker Compose ile ayağa kaldırın:

```bash
docker-compose up -d
```

**Önemli Not:** Ollama servisi çalıştıktan sonra, vektör arama işlemleri için gerekli olan embedding modelini indirmeniz gerekmektedir. Bunu sadece bir kez yapmanız yeterlidir:

```bash
docker exec -it marketplace-ollama ollama pull nomic-embed-text
```

### Servisleri Çalıştırma

Her bir mikroservisi kendi dizininden veya kök dizinden `go run` komutu ile başlatabilirsiniz.

**API Gateway'i Başlatma:**
```bash
go run cmd/api-gateway/main.go
```

**Product Service'i Başlatma:**
```bash
go run cmd/product-service/main.go
```

Benzer şekilde `user-service` ve `seller-service` de çalıştırılabilir.

## 📚 API Dokümantasyonu (Swagger)

Proje, API dokümantasyonu için Swagger kullanmaktadır.

### Erişim

API Gateway çalıştırıldıktan sonra, aşağıdaki adresten Swagger arayüzüne erişebilirsiniz:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Dokümantasyonu Güncelleme

API endpoint'lerinde değişiklik yaptığınızda veya yeni endpoint eklediğinizde, Swagger dokümantasyonunu güncellemek için proje ana dizininde aşağıdaki komutu çalıştırın. Bu komut, API Gateway ve diğer servislerin (User, Product) controller ve domain paketlerini tarar:

```bash
swag init -g cmd/api-gateway/main.go -d cmd/api-gateway,internal/api-gateway/handlers,internal/user-service/transport/http/controller,internal/user-service/domain,internal/product-service/transport/http/controller,internal/product-service/domain -o docs
```

Not: `swag` komutu yüklü değilse, aşağıdaki komutla yükleyebilirsiniz:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## 📂 Proje Yapısı

```
marketplace-backend/
├── cmd/                 # Servislerin giriş noktaları (main uygulmaları)
│   ├── api-gateway/
│   ├── product-service/
│   ├── user-service/
│   └── seller-service/
├── internal/            # Servislere özel private kodlar (business logic, repository, vb.)
├── pkg/                 # Servisler arası paylaşılan kodlar (yardımcı kütüphaneler)
└── docker-compose.yml   # Docker altyapı tanımları
```
