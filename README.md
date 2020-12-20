# GoLive

Hi all, this application features an online shop where you can sell items. 
The shop features include: 

1) a user account for each individual customer,
2) a database to store the items' information, 
3) a shopping cart to store the customer's orders, 
4) a calculator to calculate the cost of the items, and
5) security such as https, input validation, and cookie expire parameters. 

The shop can still be modified to create a platform where users can upload and sell their items. 
I have created a REST API (items.go) modifying the original code block provided by the adjunct lecturer, Wei-Meng (website: http://www.learn2develop.net/) to demonstrate my understanding of what REST means.  

With the REST API, sellers can retrieve, insert, modify or delete their item records in their respective databases. 
However, to use the file, one must have a CA signed certificate otherwise the commands will not work. 
To circumvent this, you can simply remove the cert and key files and modify ListenAndServeTLS to just ListenAndServe. 
But of course, this is only for illustration purposes as this will be a big security loophole without the SSL layer.  

The main.go file is also a modification of the original code block created by the GoSchool (https://www.goschool.sg/) lecturers, Kheng Hian (aka Ben) and Ching Yun. 
