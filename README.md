# API-Mutant (Challenge técnico)

### Verificación de ADN:

Se debe realizar un post al endpoint http://......../mutant  enviando la secuencia de ADN mediante un HTTP POST con un Json el cual tenga el
siguiente formato:

POST → /mutant/
{
“dna”:["ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACTG"]
}

En caso de verificar un mutante, devolvera un HTTP 200-OK, en caso contrario un 403-Forbidden

El metodo ademas de verificar si es mutante o no, verifica que la secuencia de ADN sea Valida y verifica si esta en db. si esta, No continua con la verifica y devuelve el resultado anterior. Si no esta en DB verifica y guarda.

### Estadisticas:

Se debe realizar un get a la siguiente URL: http://...../stats

Será devuelto un JSON con el siguiente formato de ejemplo:

{"count_mutant_dna":4,"count_human_dna":13,"ratio":0.030769}


## Instrucciones de ejecutarlo localmente:

La API esta totalmente desarrollada en GO 1.11, como base de datos se esta utilizando mongoDB.

Como intalar GO, https://golang.org/doc/install

Como instalar MongoDB, 
	Debian
	https://docs.mongodb.com/manual/tutorial/install-mongodb-on-debian/
	Ubuntu
	https://docs.mongodb.com/manual/tutorial/install-mongodb-on-ubuntu/
	Windows
	https://docs.mongodb.com/manual/tutorial/install-mongodb-on-windows/
	OSx
	https://docs.mongodb.com/manual/tutorial/install-mongodb-on-os-x/

En el proyecto se incluyeron las siguientes librerías

