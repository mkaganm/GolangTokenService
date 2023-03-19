# GolangTokenService

# Golang ile Token Servis

Örnek için hazırladığım repo ve postman collection’a bu linkten erişebilirsiniz. https://github.com/mkaganm/GolangTokenService

JSON Web Token (JWT’ler), çevrimiçi kimlik doğrulama için popüler bir yöntemdir ve JWT kimlik doğrulamasını herhangi bir sunucu tarafı programlama dilinde uygulayabilirsiniz.

Bu yazımda golang-jwt paketini kullanarak Go uygulamalarınızda JWT kullanımını anlatacağım.

## Başlarken
Go env ayarladıktan ve go.mod’u başlattıktan sonra, golang-jwt paketini kurmak için çalışma alanı dizinindeki terminalinizde şu komutu kullanabilirsiniz:
```bash
go get github.com/golang-jwt/jwt
```
## Token Oluşturma
golang-jwt paketini kullanarak JWT token oluşturmak için gizli bir anahtara ihtiyacınız olacak. Gizli anahtarınız için kriptografik olarak güvenli bir dize kullanmalı ve bunu bir ortam değişkenleri dosyasından yüklemelisiniz.

```go
func CreateToken() (string, error) {

 token := jwt.New(jwt.SigningMethodHS256)
 claims := token.Claims.(jwt.MapClaims)
 claims["exp"] = time.Now().Add(time.Hour).Unix()

 tokenStr, err := token.SignedString([]byte(models.ApiKey{}.GetXApiKey().SecretKey))
 errors.CheckErr(err)

 return tokenStr, nil
}
```
Yukarıda token oluşturmak için yazdığım metodu görebilirsiniz. JWT’i değiştirmek için Claims metodunu kullanabilirsiniz.

```go
claims["exp"] = time.Now().Add(time.Hour).Unix()
```

Burada time modülünü kullanarak tokenler için 1 saatlik bir süre belirledim.
```go
tokenStr, err := token.SignedString(models.ApiKey{}.GetXApiKey().SecretKey)
```
Bir JWT oluşturmanın son kısmı gizli anahtarı kullanmaktır. Burada SignedString metodu içinde anahtarımızı kullandık.

## Oluşturulan Tokeni Vermek İçin

```go
func GetToken(w http.ResponseWriter, r *http.Request) {

 NotPost(w, r)
 SetContentJson(w)

 if r.Header["X-Api-Key"] != nil {

  if auth.CheckPassword(r.Header["X-Api-Key"][0], []byte(models.ApiKey{}.GetXApiKey().XApiKey)) {

   token, err := CreateToken()
   errors.CheckErr(err)

   log.Default().Print("success token")

   m := make(map[string]string)
   m["token"] = token
   m["status"] = "success"

   jsonS, _ := json.Marshal(m)

   _, err = w.Write(jsonS)
   errors.CheckErr(err)

  } else {
   SendUnAuthWrite(w)
  }
 } else {
  SendUnAuthWrite(w)
 }
}
```
Yukarıda oluşturduğumuz tokeni sunan yazdığım handler fonksiyonu bulunuyor. Header içinde gelen “X-Api-Key” doğruysa tokeni kullanıcıya response içinde veriyoruz.

## JWT Token’i Doğrulama
JWT’leri doğrulamanın yöntemi ara yazılım kullanmaktır. Burada handler fonksiyonumuzu tokenimiz için oluşturulan handler fonksiyonu içinden kullanarak doğrulayabiliriz. Bir isteğin yetkili olduğunu doğrulamak için ara yazılımın nasıl kullanılacağını dah detaylıca aşağıda açıklayacağım.

```go
func ValidateToken(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//
}
```
Handler fonksiyonlarımızı yukarıda verdiğim ValidateToken fonksiyonu içinden geçirerek kullandım ve doğrulama işlemlerini bu fonksiyonun içinde gerçekleştirdim.

```go
r.HandleFunc("/token", handlers.GetToken)
r.Handle("/", handlers.ValidateToken(handlers.Index))
```
Gelen istek yetkilendirilmişse JWT fonksiyonumuz parametre olarak iletilen handler fonksiyonunu bizlere dönecek.

```go
func ValidateToken(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

  if r.Header["Authorization"] != nil {

   tokenString := r.Header["Authorization"][0]
   tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

   token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
    _, ok := t.Method.(*jwt.SigningMethodHMAC)
    if !ok {
     SendUnAuthWrite(w)
    }
    return []byte(models.ApiKey{}.GetXApiKey().SecretKey), nil
   })

   if err != nil {
    SendUnAuthWrite(w)
   }

   if token.Valid {
    next(w, r)
   }
  } else {
   log.Default().Print("Header.Token is empty")
   SendUnAuthWrite(w)
  }
 })
}
```
Yukarıda gelen isteğin tokenini kontrol eden metodu verdim. İlk olarak header’a bakacağız. Gelen isteğin header içinde ben “Authorization” olarak tokenin verilmesini bekliyorum. Aşağıda yazdığım fonksiyon tokenin yanlış verilmesi veya boş header yollanması gibi durumlarda bize UnAuthorized bir response verecek.

```go
func SendUnAuthWrite(w http.ResponseWriter) {

 SetContentJson(w)
 log.Default().Print("unauthorized trying")

 w.WriteHeader(http.StatusUnauthorized)
 m := make(map[string]string)
 m["status"] = "unauthorized"
 m["message"] = "you are unauthorized!"
 jsonM, _ := json.Marshal(m)

 _, err := w.Write(jsonM)
 errors.CheckErr(err)
}
```
Header.[“Authorization”] içinden tokenimizi alıyoruz ve aşağıdaki metod içinde kontrolünü gerçekleştiriyoruz. Bize jwt.Parse metodu bir interface ve bir hata döner. Eğer bir hata gelmişse gelen isteği UnAuthorized olarak kabul ediyoruz ve o şekilde response dönüyoruz.

```go
   token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
    _, ok := t.Method.(*jwt.SigningMethodHMAC)
    if !ok {
     SendUnAuthWrite(w)
    }
    return models.ApiKey{}.GetXApiKey().SecretKey, nil
   })
```
Yukarılarda bahsettiğim gizli anahtarı tokeni doğrulamak için kullanıyoruz ve bu metodun içinde dönüyoruz. Burada HMAC (Hash-based message authentication code) kullandım.

```go
if token.Valid {
    next(w, r)
   }
```
Yukarıda verdiğim koşul karşılandığında yani tokenimiz geçerli olduğunda yazmış olduğum ValidateToken fonksiyonun parametrelerinde gelen handler fonksiyonumuz çalışacaktır.

## Main Fonksiyonu
```go
func main() {
 r := http.NewServeMux()
 r.HandleFunc("/token", handlers.GetToken)
 r.Handle("/", handlers.ValidateToken(handlers.Index))

 log.Default().Print("Server localhost:8081 is started")

 defer func(addr string, handler http.Handler) {
  err := http.ListenAndServe(addr, handler)
  errors.CheckErr(err)
 }(":8081", r)
}
```
Main fonksiyonunun içinde bir HTTP server ayaklandırdık ve bunu 8081 portundan vererek gelen istekleri dinlemeye başladık. Yukarıda verdiğim gibi handler fonksiyonlarımızı oluşturduğum ValidateToken fonksiyonu içinde kullanarak tokenimizi doğrulamış olduk.

Postman ile Çalıştırılması
Postman üzerinde nasıl çalıştığına hep birlikte bir göz atalım. Kullandığım postman koleksiyonuna repo içinden ulaşabilirsiniz. https://github.com/mkaganm/GolangTokenService


Bu şekilde bir postman environment’i oluşturdum.


/token endpoint e post isteği attığımızda bizlere response içinde tokeni veriyor.


Postman içinde Tests kısmına yazdığımız ufak bir JavaScript kodu ile gelen tokeni environment içine kaydediyoruz.

```js
var res = pm.response.json();
pm.environment.set('Bearer', res.token);
```
Environment içine kaydettiğimiz tokeni Bearer token kısmına verdiğimizde bizlere başarılı bir dönüş sağlıyor ve HTTP 200 dönüyor.


Yanlış bir token verdiğimizde ise bize UnAuthorized mesajı ve HTTP 401 dönüyor.

Umarım faydalı olmuştur. Okuduğunuz için teşekkür ederim. Başka yazılarda görüşmek üzere hoşça kalın.
