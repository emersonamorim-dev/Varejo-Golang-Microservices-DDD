# Varejo-Golang-Microservices-DDD

Codificação em Golang com uso de Framework GIN o projeto **Varejo-Golang-Microservices** é uma solução completa e robusta para gerenciamento de varejo, implementado usando Go (Golang) e modelado de acordo com o padrão de modelagem Domain-Driven Design (DDD). Com uma arquitetura de microserviços, essa solução abrange várias funcionalidades do setor de varejo, desde gestão de clientes até relatórios e suporte.

## Características Principais

- **DDD (Domain-Driven Design)**: A solução segue estritamente os princípios do DDD, garantindo um design orientado a domínio, o que facilita a modelagem de acordo com as regras de negócios específicas do setor de varejo.

- **Banco de Dados MongoDB**: Utilizado para armazenamento persistente, aproveitando sua escalabilidade e flexibilidade para lidar com grandes volumes de dados.

- **Kafka**: Adotado para mensageria assíncrona entre os microserviços, garantindo comunicação eficaz e resiliente.

- **Docker**: Cada microserviço é containerizado usando o Docker, assegurando isolamento, reprodutibilidade e escalabilidade.

- **Kubernetes**: Utilizado para orquestração dos contêineres, proporcionando gerenciamento eficiente, balanceamento de carga e autorecuperação.

- **Terraform**: Para provisionamento e gestão da infraestrutura como código.

- **Jenkins**: Integração contínua e entrega contínua (CI/CD) para automatizar o processo de build e deploy.

## Arquitetura e Padrões

- **Design Patterns**: A solução utiliza diversos padrões de design para resolver problemas comuns, incluindo Factory, Singleton e Strategy, promovendo código reutilizável e manutenível.

- **SOLID Principles**: O código segue os princípios SOLID, garantindo um design de software robusto, escalável e fácil de manter.

- **Programação**: A aplicação é predominantemente orientada a objetos, aproveitando as vantagens de encapsulamento, herança e polimorfismo.



## Microserviços Varejo

#### Rode a aplicação de forma local com comando:

```
go run main.go
```

#### Para rodar cada Microservice local em separado acesse o diretório:
```
/customer-service/api
```
#### Isso para acessa cada Microservice em separado
```
go run main.go
```


### 1. `customer-service`
Responsável pela gestão de clientes, incluindo criação, atualização, recuperação e exclusão de informações do cliente.

### 2. `integration-service`
Facilita a integração com sistemas e serviços externos, atuando como uma ponte entre a plataforma e terceiros.

### 3. `location-service`
Gerencia informações relacionadas a localizações, como endereços de entrega, lojas físicas e centros de distribuição.

### 4. `order-service`
Encarregado do ciclo de vida dos pedidos, desde a criação até a conclusão, incluindo acompanhamento e status.

### 5. `payment-service`
Trata de todas as transações financeiras, processamento de pagamentos, faturas e integração com gateways de pagamento.

### 6. `product-service`
Gerencia o catálogo de produtos, incluindo detalhes, preços, inventário e categorizações.

### 7. `promotion-service`
Responsável por campanhas promocionais, descontos e cupons.

### 8. `report-service`
Fornece relatórios detalhados sobre vendas, clientes, produtos e outras métricas vitais.

### 9. `support-service`
Oferece suporte ao cliente e gerencia tickets, reclamações e feedback.

### Diagrana da Aplicação:



 # 🛠 Tecnologias Usadas

## ⚙️ Backend

- **DDD (Domain-Driven Design)**
  
- **Banco de Dados MongoDB**
  
- **Kafka**
  
- **Docker**
  
- **Kubernetes**
  
- **Terraform**

- **Jenkins**




## Conclusão

**Varejo-Golang-Microservices** é uma solução de ponta projetada para atender às necessidades de negócios complexos no espaço de varejo. Com uma combinação de tecnologias modernas e padrões de design, oferece escalabilidade, confiabilidade e facilidade de manutenção.




### Autor:
Emerson Amorim
