import Interfaces.IDataBase;
import Interfaces.ITable;
import main.DataBase;
import main.Table;

import java.util.List;

public class Main {
    public static void main(String[] args) {
        IDataBase dataBase = new DataBase();
        ITable table = new Table();
        table.add("workKey", "workValue", 0);
        dataBase.add("dataBaseKey", table);
        System.out.println(dataBase.get("dataBaseKey").get("workKey"));
    }
}