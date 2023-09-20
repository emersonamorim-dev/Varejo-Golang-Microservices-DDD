# Varejo-Golang-Microservices-DDD

## Introdução

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

## Microserviços

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

+-----------------------------+
|     Varejo-Golang-          |
|     Microservices           |
+-----------------------------+
               |
   +-----------v------------+
   |       API Gateway      |
   +-----------+------------+
               |
 +-------------v----------+    +---------------------+
 |  customer-service      |<-->|   MongoDB Database  |
 +-----------------------+    +---------------------+
               |
 +-------------v----------+    +---------------------+
 |  location-service      |<-->|   Kafka Messaging   |
 +-----------------------+    +---------------------+
               |
 +-------------v----------+    +---------------------+
 |  integration-service   |<-->|     Docker Engine   |
 +-----------------------+    +---------------------+
               |
 +-------------v----------+    +---------------------+
 |  order-service         |<-->|   Kubernetes Cluster|
 +-----------------------+    +---------------------+
               |
 +-------------v----------+    +---------------------+
 |  payment-service       |<-->|    Jenkins CI/CD    |
 +-----------------------+    +---------------------+
               |
 +-------------v----------+    +---------------------+
 |  product-service       |<-->|    Terraform IaC    |
 +-----------------------+    +---------------------+
               |
 +-------------v----------+
 |  promotion-service     |
 +-----------------------+
               |
 +-------------v----------+
 |  report-service        |
 +-----------------------+
               |
 +-------------v----------+
 |  support-service       |
 +-----------------------+


 # 🛠 Tecnologias Usadas

## ⚙️ Backend

- ![DDD Icon](URL_TO_ICON_IMAGE_FOR_DDD) **DDD (Domain-Driven Design)**
  
- ![MongoDB Icon](https://camo.githubusercontent.com/9ebde7ca22ab3f3b4bf92d2743804ab9e581e413a16cdf3626c2092e69967d80/68747470733a2f2f63646e2e6a7364656c6976722e6e65742f67682f64657669636f6e732f64657669636f6e2f69636f6e732f6d6f6e676f64622f6d6f6e676f64622d6f726967696e616c2e737667) **Banco de Dados MongoDB**
  
- ![Kafka Icon](https://github.com/devicons/devicon/blob/v2.15.1/icons/apachekafka/apachekafka-original-wordmark.svg) **Kafka**
  
- ![Docker Icon](https://camo.githubusercontent.com/f64a041d6d0cda76988a117724ce3b3272b8fc5f9f742c4dcb9160be9a2c41c1/68747470733a2f2f63646e2e6a7364656c6976722e6e65742f67682f64657669636f6e732f64657669636f6e2f69636f6e732f646f636b65722f646f636b65722d706c61696e2e737667) **Docker**
  
- ![Kubernetes Icon](https://github.com/devicons/devicon/blob/v2.15.1/icons/kubernetes/kubernetes-plain-wordmark.svg) **Kubernetes**
  
- ![Terraform Icon](URL_TO_ICON_IMAGE_FOR_TERRAFORM) **Terraform**

- ![Jenkins Icon](https://raw.githubusercontent.com/devicons/devicon/v2.15.1/icons/jenkins/jenkins-original.svg) **Jenkins**




## Conclusão

**Varejo-Golang-Microservices** é uma solução de ponta projetada para atender às necessidades de negócios complexos no espaço de varejo. Com uma combinação de tecnologias modernas e padrões de design, oferece escalabilidade, confiabilidade e facilidade de manutenção.




### Autor:
Emerson Amorim
