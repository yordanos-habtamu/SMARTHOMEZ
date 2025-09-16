

export {}
declare global{
    type Property ={
        id:string
        title:String
        description:string
        image:string
        category:string
        bedrooms:number
        bathrooms:number
        squareFeet:number
        price:number
        listedDate:string
    }
    type Category ={
        id:string
        name:string
        descritpion:string
        image:string
    }
}