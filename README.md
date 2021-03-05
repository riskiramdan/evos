# Evos Test Assesement using Golang by Riski Ramdan

## Getting started (< 2mn)

```
git clone git@github.com:riskiramdan/evos.git
cd evos

docker-compose up //Running PostgresSQL & Golang Application
```
## Running Go Application 
```
docker-compose up
```

## Postman Documentation
https://documenter.getpostman.com/view/9740098/Tz5jeLBV
https://www.getpostman.com/collections/5980f656d7d002e04fb6

## Domain Driven Design Architectures

Software design is a very hard thing. From years, a trend has appeared to put the business logic, a.k.a. the (Business) Domain, and with it the User, in the heart of the overall system. Based on this concept, different architectural patterns was imaginated. 

One of the first and main ones was introduced by E. Evans in its [Domain Driven Design approach](http://dddsample.sourceforge.net/architecture.html).

![DDD Architecture](/doc/DDD_architecture.jpg)

Based on it or in the same time, other applicative architectures appeared like [Onion Architecture](https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/) (by. J. Palermo), [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) (by A. Cockburn) or [Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) (by. R. Martin).

This repository is an exploration of this type of architecture, mainly based on DDD Architecture, on a concrete and modern golang application.