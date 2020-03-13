pragma solidity 0.6.4;
pragma experimental ABIEncoderV2;

import "./library/SafeMath.sol";
import "./library/ERC20.sol";
import "./library/Ownable.sol";
import "./library/ERC20Detailed.sol";

struct Beverage {
    string name;
    uint16 price;
    uint8 amount;
}

contract Token is ERC20, ERC20Detailed, Ownable{
    uint8 private _maxAmount = 50;
    uint8 private _maxKinds = 20;
    uint16 private _profit = 200;
    Beverage[] allBeverages;
    
    event AddBeverage(uint8 indexed index, string name, uint16 price, uint8 amount);
    event RemoveBeverage(uint8 indexed index, string name, uint16 price, uint8 amount);
    event BuyBeverage(uint8 indexed index, string name, uint16 price);
    constructor()
    ERC20Detailed("TempToken","Temp",18) public {
        super._mint(_msgSender(), 500000000*10**18);
        super._mint(address(this), 500000000*10**18);
    }

    function managedBalace() public view returns(uint256){
        return balanceOf(address(this));
    }

    function addBeverage(string memory _name, uint16 _price, uint8 _amount) public onlyOwner {
        require(allBeverages.length < _maxKinds,"already Full Kinds");
        super._burn(address(this), (_price-_profit)*_amount);
        allBeverages.push(Beverage({
            name : _name,
            price: _price,
            amount: _amount
        }));
        uint8 index = uint8(allBeverages.length.sub(1));
        emit AddBeverage(index, _name, _price, _amount);
    }

    function removeBeverage(uint8 _index) public onlyOwner {
        require(allBeverages.length > _index,"Does not exist");
        Beverage memory removedB = allBeverages[_index];
        if(_index != allBeverages.length-1) {
            Beverage memory b = allBeverages[allBeverages.length-1];
            allBeverages[_index] = b;
        }
        allBeverages.pop();
        uint totalPrice = (removedB.price - _profit) *removedB.amount;
        super._mint(address(this), totalPrice);
        emit RemoveBeverage(_index, removedB.name, removedB.price, removedB.amount);
    }

    function fillMaxAmount(uint8 _index) public onlyOwner {
        Beverage storage b = allBeverages[_index];
        require(b.amount < _maxAmount,"already Full Amount");
        uint amount = _maxAmount - b.amount;
        uint totalPrice = (b.price - _profit)*amount;
        
        super._burn(address(this), totalPrice);
        b.amount = _maxAmount;
    }

    function allBeveragesLength() public view returns(uint8) {
        return uint8(allBeverages.length);
    }
    
    function showBeverageByIndex(uint8 _index) public view returns (uint8,string memory,uint16,uint8) {
        return (_index,allBeverages[_index].name,allBeverages[_index].price,allBeverages[_index].amount);
    }

    function showBeverages() public view returns
    (
        uint8[] memory index, 
        string[] memory names, 
        uint16[] memory prices, 
        uint8[] memory amount
    ) 
    {
        index = new uint8[](allBeverages.length);
        names = new string[](allBeverages.length);
        prices = new uint16[](allBeverages.length);
        amount = new uint8[](allBeverages.length);
        for(uint8 i = 0 ; i < allBeverages.length ; i++) {
            index[i] = i;
            names[i] = allBeverages[i].name;
            prices[i] = allBeverages[i].price;
            amount[i] = allBeverages[i].amount;
        }
    }

    function buyBeverage(uint8 _index) public {
        Beverage storage b = allBeverages[_index];
        require(b.amount > 0,"sold out");
        super.transfer(address(this), b.price);
        emit BuyBeverage(_index, b.name, b.price);
        b.amount = b.amount - 1;
    }
}