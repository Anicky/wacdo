```mermaid
classDiagram

    class Item {
	    +int id
	    +string name
	    +string description
	    +string image
	    +float price
	    +bool isAvailable
	    +datetime createdAt
	    +datetime updatedAt
    }

    class Product {
	    +int id
    }

    class ProductCategory {
	    +int id
	    +string name
	    +string description
    }

    class Menu {
	    +int id
    }

    class MenuOption {
	    +int id
    }

    class Order {
	    +int id
	    +string ticketNumber
	    +datetime createdAt
	    +datetime preparedAt
	    +datetime deliveratedAt
	    +getPrice()
	    +setStatus(status)
	    +isReadyForDelivery()
    }

    class User {
	    +int id
	    +string username
	    +string password
	    +bool isActive
	    +datetime createdAt
	    +datetime updatedAt
	    +checkPassword(password)
    }

    class UserRole {
	    ADMIN
	    GREETER
	    ORDER_PICKER
	    MANAGER
    }

    class OrderStatus {
	    CREATED
	    IN_PREPARATION
	    PREPARED
	    DELIVERED
    }

	<<enumeration>> UserRole
	<<enumeration>> OrderStatus

    Product "*" --> "1" ProductCategory : has
    User "*" <--> "*" UserRole : has
    Order "*" --> "1" User : isManagedBy
    Order "*" --> "1" OrderStatus : has
    Product --|> Item
    Menu --|> Item
    Menu "*" --> "*" Product : contains
    Order "*" --> "*" Item : contains
    MenuOption --|> Item
    Menu "*" -- "*" MenuOption : contains
```