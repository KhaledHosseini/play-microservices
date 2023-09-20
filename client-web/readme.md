Contents:
- [Summary](#summary)
- [Tools](#tools)
- [Docker dev environment](#docker-dev-environment)
- [API mock service: Wiremock](#api-mock-service-wiremock)
- [Client-web service: Typescript](#client-web-service-typescript)
- [To DO](#to-do)


---

This is the 7th part of a series of articles under the name **"Play Microservices"**. Links to other parts:<br/>
Part 1: [Play Microservices: Bird's eye view](https://dev.to/khaledhosseini/play-microservices-birds-eye-view-3d44)<br/>
Part 2: [Play Microservices: Authentication](https://dev.to/khaledhosseini/play-microservices-authentication-4di3)<br/>
Part 3: [Play Microservices: Scheduler service](https://dev.to/khaledhosseini/play-microservices-scheduler-19km)<br/>
Part 4: [Play Microservices: Email service](https://dev.to/khaledhosseini/play-microservices-email-service-1kmc)<br/>
Part 5: [Play Microservices: Report service](https://dev.to/khaledhosseini/play-microservices-report-service-4jcm)<br/>
Part 6: [Play Microservices: Api-gateway service](https://dev.to/khaledhosseini/play-microservices-api-gateway-service-4a9j)<br/>
Part 7: You are here<br/>
Part 8: [Play Microservices: Integration via docker-compose](https://dev.to/khaledhosseini/play-microservices-integration-via-docker-compose-2ddc)<br/>
Part 9: [Play Microservices: Security](https://dev.to/khaledhosseini/play-microservices-security-45e4)<br/>

---

## Summary

In the previous stages, we successfully developed all of back-end services including the auth, scheduler, email, report and api-gateway services. Our current objective is to establish a client service that acts as a single entry point for end users to access our application. As we independently develop the client service, the remaining services are unavailable to us during development. To overcome this limitation, we create mock implementations to simulate the behavior of the unavailable services during the development of the client service.

![dev environment](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/gg6xpnbh1wcv2p186myf.png)


At the end, the project directory structure will appear as follows:


![Folder structure](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/p0v2u3tnz7vh0lwqrd6d.PNG)

---


## Tools

The tools required In the host machine:

  - [Docker](https://www.docker.com/): Containerization tool
  - [VSCode](https://code.visualstudio.com/): Code editing tool
  - [Dev containsers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension for VSCode
  - [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) extension for VSCode
  - [Git](https://git-scm.com/)

The tools and technologies that we will use Inside containers for each service:

 - Wiremock service: [Wiremock](https://hub.docker.com/r/wiremock/wiremock). We use this service to simulate our api-gateway service.
 - Web-Client service: 
  - [Typescript](https://www.typescriptlang.org/) : programming language
  - [Next.js](https://nextjs.org/): Web framework that can be used for developing backend and front-end applications.
  - [Tailwind css](https://tailwindcss.com/):  A utility-first CSS framework 
  - [react-hook-form](https://react-hook-form.com/docs)
  - [TanStack Query](https://tanstack.com/query/latest): A asynchronous state management.

---

## Docker dev environment

Development inside Docker containers can provide several benefits such as consistent environments, isolated dependencies, and improved collaboration. By using Docker, development workflows can be containerized and shared with team members, allowing for consistent deployments across different machines and platforms. Developers can easily switch between different versions of dependencies and libraries without worrying about conflicts.

![dev container](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/id0cijn8hm77ohn0hk2j.png)

When developing inside a Docker container, you only need to install `Docker`, `Visual Studio Code`, and the `Dev Containers` and `Docker `extensions on VS Code. Then you can run a container using Docker and map a host folder to a folder inside the container, then attach VSCode to the running container and start coding, and all changes will be reflected in the host folder. If you remove the images and containers, you can easily start again by recreating the container using the Dockerfile and copying the contents from the host folder to the container folder. However, it's important to note that in this case, any tools required inside the container will need to be downloaded again. Under the hood, When attaching VSCode to a running container, Visual Studio code install and run a special server inside the container which handle the sync of changes between the container and the host machine.

---

## API mock service: Wiremock

During the development of our microservice application, we have implemented various patterns to ensure efficient development. One of the patterns we have followed is service-per-team development. This approach focuses on each team developing their services independently, with limited knowledge of and no direct access to other services. When our service relies on another service that is inaccessible during development, we use mocking techniques to simulate the behavior of that service. For different service communication protocols such as gRPC, REST API, GraphQL, and others, we have various applications and even online services that do this job for us. In our current scenario, we aim to replicate the behavior of an API gateway service. To achieve this, we rely on the usage of [Wiremock](https://wiremock.org/), a tool that offers a convenient solution for mocking APIs. By utilizing Wiremock, we can simulate the responses and behavior of the API gateway service, allowing us to continue development and testing seamlessly. Using Wiremock and similar tools, we can effectively emulate the behavior of external services, enabling smoother development and testing workflows within our microservice architecture. Running and configuring Wiremock via Docker is quite easy. Lets begin!

> - Create a folder for the project and choose a name for it (such as 'microservice'). Then create a folder named `client-web`. This folder is the root directory of the current project. You can then open the root folder in VS Code by right-clicking on the folder and selecting 'Open with Code'.
> - Inside the root directory create a folder with the name `wiremock`, then create a Dockerfile and set content to `FROM wiremock/wiremock:2.35.0`
> - Create a folder inside wiremock named mappings. Inside this folder we define our endpoints. For example lets say we have the following endpoint in our api-gateway:


![End point example](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/f5pgvmo1waphb4y308mj.PNG)
Then using a json file containing a request and a corresponding response we can configure wiremock to mimic that end-point.

```json
{
    "request": {
      "method": "GET",
      "urlPattern": "/api/v1/ping"
    },
    "response": {
      "status": 200,
      "body": "{\"Message\": \"Pong\"}"
    }
}
```
> - During the development of the client service, our API gateway service is not accessible. However, the protocol layer models for each service have already been determined and made accessible through collaborative efforts and guidance from the Technical Leads. We use this models to mock our api-gateway service. Copy all json files from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/client-web/wiremock/mappings) to mapping folder.
> - Create a file named .env in the root directory and add the following content:

```
WIREMOCK_PORT = 8088
WIREMOCK_CONTAINER_PORT=8080
CLIENT_PORT = 3000
```
> - Inside root directory create a file named docker-compose.yml and add the following content.

```yaml
version: '3'
services:
  wiremock:
    build: 
      context: ./wiremock
      dockerfile: Dockerfile
    container_name: wriremock
    ports:
      - ${WIREMOCK_PORT}:${WIREMOCK_CONTAINER_PORT}
    volumes:
      - ./wiremock/mappings:/home/wiremock/mappings
```

> - Run `docker-compose up -d --build`. Now go to `http://localhost:8088/__admin/`. You can see the available mock end points. 


![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/ciqepbk9m7fidjeuoqga.PNG)

> - you can test this service using applications like [postman](https://www.postman.com/).
> - Our api-gateway mock is ready! run `docker-compose down`

---

## Client-web service: Typescript

Before delving into the specifics, let's begin by describing the application we are developing. The application consists of four backend services: authentication, job scheduler, and report services. These services are not directly accessed by the client application. Instead, an API gateway acts as a bridge, facilitating communication between the client and these services.
The client application offers the following functionalities:
1. User Registration and Login: The application provides signup and login pages where users can create and authenticate their accounts.
2. Admin Features: Administrative users have additional capabilities, including:
   - Querying the List of Registered Users: Admins can retrieve information about the registered users.
   - Querying the List of Scheduled Jobs: Admins have the ability to inquire about the jobs currently scheduled and edit or delete them.
   - Scheduling New Jobs: Admins can create and schedule new jobs.
   - Querying Reports: Admins can retrieve reports generated by the system.

![client summary](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/i1mc9izdoeb6pa0v6xiu.png)

> - Our Next.js application structure can be summarized as follow: 
   1. A server that acts as a proxy between the api gateway and the client app. This server is optional and one benefit of it is Solving [CORS ](https://en.wikipedia.org/wiki/Cross-origin_resource_sharing) problems as our api-gateway and web-client service may run on different domains. Also we can add an extra layer of security in this server to protect our-api gateway.
   2. Three main pages that all are rendered in the [client side](https://nextjs.org/learn/foundations/how-nextjs-works/rendering) (inside browser). Signup, login and main page. the main page has three components. One for querying the users list, One for handling jobs and the final one for querying reports.

> - Create a folder named `client-service` inside `client-web` folder.
> - Create a Dockerfile inside `client-service` and set the contents to

```shell
FROM node:20.4.0

WORKDIR /usr/src/app
```
> - Add the following to the service part of our docker-compose.yml file.

```yaml
    client:
    build: 
      context: ./client-service
      dockerfile: Dockerfile
    container_name: client
    command: sleep infinity
    ports:
      - ${CLIENT_PORT}:${CLIENT_PORT}
    environment:
      - APIGATEWAY_URL=http://wriremock:${WIREMOCK_CONTAINER_PORT}
    volumes:
      - ./client-service:/usr/src/app
```
> - We are going to do all the development inside a docker container without installing node.js in our host machine. To do so, we run the containers and then attach VSCode to the client-service container. As you may noticed, the Dockerfile for client-service has no entry-point therefore we set the command value of it to `sleep infinity` to keep the container awake.
> - Now run `docker-compose up -d --build`
> - While running, attach to the client service by clicking bottom-left icon and then select `attach to running container`. Select client service and wait for a new instance of VSCode to start. Upon starting the attached instance of VSCode, you will be prompted to open a folder within the container. As we have designated the WORKDIR as /usr/src/app inside the Dockerfile, we will select this folder inside the container. It is important to note that this designated folder is mounted to the client-service folder on the host machine using Docker Compose volumes. Consequently, any changes made within the selected folder will be automatically synced to this folder on the host machine. This synchronization ensures that modifications made during development are reflected in both the container and the host environment.
> - After opening the folder `/usr/src/app`, open a new terminal and initialize the next.js project by running `npx create-next-app --typescript client`. You need to go through a list of question before initializing the project. Select the answers as shown below.


![create project questions](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/8eoc7oqdx621tui7o0bl.PNG)

> - After initializing the app, a file named packages.json is created. This file contains information such as dependencies, package name, version, etc. Another part of this file is the scripts. 

```json
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
```
> - We can run this scripts from the client folder using `npm run <scriptname>`. Run `cd client` and then run `npm run dev`. Go to `http://localhost:3000` and you can see the default page of our project. stop the service by hitting `ctl + c`
> - We use newly introduced `app routing` in our project. To compare app-routing and page routing refer to [next.js](https://nextjs.org/docs) website and read the documentation by selecting them.

![Next.js routing](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/72gtv055kne970yza3v5.PNG)
> - We start by creating our server. This server will act as a proxy between our website and the api gateway. Create a folder named config inside src folder and then a file named index.tsx. set the content to

```typescript
export const URL_APIGATEWAY=`${process.env.APIGATEWAY_URL}/api/v1`;
```
> - In the app router routing method, You can create a folder XXX inside the app folder and then create a file named rout.tsx and a file named page.tsx inside. The route.tsx acts as your api end point and the page.tsx is a webpage at that route. Create a folder named api and then a folder inside named ping. create a file named route.tsx inside ping and set the content to:

```typescript
import { NextRequest, NextResponse } from "next/server";

export async function GET(req: NextRequest) {
    return new NextResponse(JSON.stringify({ message: "Pong" }), {
        status: 200,
      })
}
```
> - Run `npm run dev`. Now go `http://localhost:3000/api/ping`. you can see the response saying `pong`!
> - Stop the server by hitting `ctl + c`
> - Create a file named middleware.tsx inside src folder and set the content to

```typescript
import { NextRequest, NextResponse } from 'next/server'
import {URL_APIGATEWAY} from './config'

export async function middleware(req: NextRequest) {
  const regex = new RegExp('/api/*')
  if (!regex.test(req.url)) {
    return new NextResponse(null, {
      status: 400,
      statusText: "Bad Request"
    })
  }
  console.log("middleware is called for url: ",req.url)
  const url = URL_APIGATEWAY + "/" + req.url.split("/api/")[1]
  console.log("middleware sends the request to : ", url)
  const res = await fetch( url, {
    method:req.method,
    headers: req.headers,
    body: req.body
  });

  return res
}

export const config = {
  matcher: '/api/:path*',
}
```
> - This middleware runs only for requests to the '/api/:path*' path. When we request this api, it  will get the response from the api-gateway.
> - Run `npm run dev`. Now go `http://localhost:3000/api/ping`. you can see the response saying `message: Pong.`. This time the result has been returned from the api-gateway (Our mock service).
> - Stop the service by hitting `ctl + c`
> - Create a folder named types inside src. Inside this folder we will define our models. Copy all the files from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/client-web/client-service/client/src/types) to this folder.
> - run `npm install @tanstack/react-query react-hook-form react-hot-toast`
> - Create a folder named components and then a folder named providers inside it. Here we are going to put the components that provide a capability for a [context](https://react.dev/reference/react/useContext). Copy the files from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/client-web/client-service/client/src/components/providers). One of the providers is AuthProvider. This component provides a context for the logged-in user to be used by any child elements inside the app. 

```typescript
'use client';

import {User} from '@/types'
import React from 'react';

const UserContext = React.createContext<
  [User | null, React.Dispatch<React.SetStateAction<User | null>>] | undefined
>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = React.useState<User | null>(null);
  return (
    <UserContext.Provider value={[user, setUser]}>
      {children}
    </UserContext.Provider>
  );
}

export function useAuth() {
  const context = React.useContext(UserContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within a UserContext');
  }
  return context;
}
```
> - We warp the whole pages in our application inside these providers. go to layout.tsx inside the app folder and change the file content to:

```typescript
import './globals.css';
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { QueryProvider } from '@/components/providers/query_provider'
import {AuthProvider} from '@/components/providers/auth_provider'
import ToastProvider from '@/components/providers/toast_provider'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Microservices job scheduler',
  description: 'A simple job scheduler app with microservices architecture.',
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {

  return (
    <html lang="en">
      <body className={`${inter.className}  font-inter antialiased bg-gray-300 text-gray-900 tracking-tight`}>
      <ToastProvider />
        <AuthProvider>
          <QueryProvider>
            {children}
          </QueryProvider>
        </AuthProvider>
      </body>
    </html>
  )
}
```

> - Create a folder named client inside components folder and then a file named signup_form.tsx. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/client-web/client-service/client/src/components/client/signup_form.tsx). Inside this file we use [`useForm`](https://react-hook-form.com/docs/useform) for binding form content to `CreateUserRequest `model and [`useMutation`](https://tanstack.com/query/latest/docs/react/guides/mutations) to perform user registration.

> - Create a folder named `(auth)` (with parenthesis: this name will be ignored in the routing) and then a folder named signup. then create a file named page.tsx and set the content to

```typescript
'use client'

import { useAuth } from "@/components/providers/auth_provider";
import SignUpForm from "@/components/client/signup_form";
import { redirect } from 'next/navigation'

export default function SignUp() {
  
  const [user,_setUser] = useAuth()
  if (user) {
      redirect("/")
  }
  
  return (
    <main>
      <SignUpForm/>
    </main>
  );
}
```
> - here we have used `useAuth()` to check if the user is already logged in or not. If yes, simply redirect to the homepage. Now run `npm run dev` and then go to `http://localhost:3000/signup`. Fill in the form to match the mapping of our mock api (In this case name: admin,email:admin@admin.com,password: password,role: 0). Now hit signup. If everything goes according to plan, the sign up would be successful and you will be redirected to the login page which does not exist at the moment.

![Signup form](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/b9i3uf32wvzph6kh3qmq.PNG)

> - Stop the server by hitting `ctl + c`
> - Create a file named login_form.tsx inside components/client folder and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/client-web/client-service/client/src/components/client/login_form.tsx).
> - Create a folder inside `(auth)` folder named login and then a file named page.tsx. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/client-web/client-service/client/src/app/(auth)/login/page.tsx).Now run `npm run dev` and then go to `http://localhost:3000/login`. Fill in the form to match the mapping of our mock api (in this case: email: admin@admin.com,password: password). If everything goes according to plan, the login would be successful and you will be redirected to the homepage.
> - Stop the server by hitting `ctl + c`

> - Create a folder named logout inside (auth) folder. Then a file named page.tsx. set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/client-web/client-service/client/src/app/(auth)/logout/page.tsx).
> Create a folder named (default) inside the app folder and then move the page.tsx file from app folder to (default). This file would be our homepage. We first create necessary components and then show them inside our home page.
> - Create a file named header.tsx inside components/client. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/client-web/client-service/client/src/components/client/header.tsx). Inside this component, we will show some items based on the logging state of the user. If the user is logged in, we show user's email and logout button. if the user is not logged in, we show login and signup buttons.
> - Create a folder named lib inside src folder. Then a file named api_gateway.tsx. We define a function called `fetch_with_refresh_token`. This function checks the return code for calls to our protected end-points and if the result is unauthorized, then we refresh the access token. Set the contents to:

```typescript
import { StatusCodes } from 'http-status-codes';
// we do fetch, If the result is unauthorized, Possibly our access token has expired! we simply refresh it.
export async function fetch_with_refresh_token(url: string, options?: RequestInit | undefined): Promise<Response> {
    let result = await fetch(url,options)
    console.log("fetch result status code and texts are: ",result.status, result.statusText)
    if (result.status == StatusCodes.UNAUTHORIZED) {
        //possibly expired access token.
        console.log("Unauthorized. possibly access token has expired. lets refresh access token....")
        result = await fetch('/api/user/refresh_token', {
            method: "POST"
        })
        console.log("refresh token result arrived",result)
        if (result.status == StatusCodes.OK) {
            console.log("refresh token request is successful. calling again.")
            //we have accessed new access token in browser cookie. now request again with our new access token
            result = await fetch(url,options)
            return result
        }else {
            console.log("refresh token request was not successful")
        }
    }
    return result
}
```
> - Change the contents of page.tsx inside (default) folder to

```typescript
'use client'

import React, { useEffect, useState } from "react"
import { ParseUser } from "@/types";
import { useAuth } from "@/components/providers/auth_provider";
import Header from '@/components/client/header'
import {fetch_with_refresh_token} from '../../lib/api_gateway'

export default function Home() {

  const [_user,setUser] = useAuth()

  useEffect(() => {
    // declare the data fetching function
    const fetchData = async () => {
      const userResponse = await fetch_with_refresh_token("/api/user/get");
      const user = await ParseUser(userResponse)
      console.log("Home.useEffect.fetchData: User is: ",user)
      setUser(user)
    }
  
    fetchData()
    .catch(console.error);
  }, [setUser])

  const [currentComponent, setCurrentComponent] = useState<string | null>(null);

  return (
    <div>
      <Header/>
      <div className="flex min-h-screen flex-row bg-gray-100 text-gray-800">
        <main className="main -ml-48 flex flex-grow flex-col p-4 transition-all duration-150 ease-in md:ml-0">
          Main
        </main>
      </div>
    </div>
  );
}
```
> - Now run `npm run dev`. Then go to `http://localhost:3000/` and you can see the header. You can now navigate through the pages. Go to `http://localhost:3000/signup`. Enter credentials and you will be redirected to login page. Then fill in the form and after hitting login button you will be redirected to the home page. As you are logged in now, the email and the logout button will be shown on the header.

![header](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/6ojngu4wz9ey9319m8zr.PNG)

> - Now it is time to create other components of our home page. Copy the remaining files from [here ](https://github.com/KhaledHosseini/play-microservices/tree/master/client-web/client-service/client/src/components/client) to components/client folder. This files are components for user, job and report handlings. Set the content of page.tsx inside (default) folder from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/client-web/client-service/client/src/app/(default)/page.tsx). 
> - Run `npm run dev`. Go to `http://localhost:3000/` and voila. Our web app is ready.


![Home page](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/hqz88x5c8l4ccekgj76z.PNG)


---


## To DO

> - Add tests
> - Add tracing using Jaeger
> - Add monitoring and analysis using grafana
> - Refactoring

---

I would love to hear your thoughts. Please comment your opinions. If you found this helpful, let's stay connected on Twitter! [khaled11_33](https://twitter.com/khaled11_33).