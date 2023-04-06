## Template description

This is a Shopify app template written in Golang, that bootstraps the app building process. It includes:

- The setup of the client and server parts, built on top of the App Bridge.
- Shopify app installation logic.
- Examples of using the Shopify API include creating store products and counting the number of products.

## Template usage

### Prerequisites

Please ensure that the following software is installed on your computer:

- [Node.js](https://nodejs.org/)
- [Golang](https://go.dev/)
- [Docker](https://www.docker.com/)

### Getting started

1. Clone the template using the following terminal command:
    
    ```
    npx @softcerycom/shopify-app-template-go-init@latest && cd softcery-shopify-app-template-go
    ```
    
2. Install NPM dependencies:
    
    ```
    npm i    
    ```
    
3. Start the database (Postgres is used by default). The following command will start a new Postgres container using the configuration in the **`.local.env`** file:
    
    ```
    docker-compose --env-file .local.env up --build postgresdb
    ```
    
4. Run the project:
    
    ```
    npm run dev
    ```
