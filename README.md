# **Template Usage Guide**

This guide will walk you through the steps required to start using the Softcery Shopify App Template.

## **Prerequisites**

Before you can use the template, you will need the following software installed on your computer:

- Node.js
- Golang
- Docker

## **Getting Started**

To get started with the template, follow these steps:

### **1. Clone the Template**

To clone the template, run the following command in your terminal:

```
npx @softcerycom/shopify-app-template-go-init@latest
```

This will create a new directory with the template files in the current directory.

### **2. Start the Database**

The Softcery Shopify App Template Go Init uses Postgres for data storage. To start the database, navigate to the root directory of the template and run the following command:

```
docker-compose --env-file .local.env up --build postgresdb
```

This will start a new Postgres container using the configuration in the **`.local.env`** file.

### **3. Start the Template**

To start the template, navigate to the root directory of the template and run the following command:

```
npm run dev
```

## **Conclusion**

You should now have a working instance of the Shopify App in Go. You can start modifying the template to build your own Shopify app!
