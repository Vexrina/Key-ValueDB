import main.DataBase;
import main.Work;

public class Main {
    public static void main(String[] args) {
        DataBase dataBase = new DataBase();
        Work work = new Work();
        work.add("workKey", "workValue", 0);
        dataBase.add("dataBaseKey", work);

        System.out.println(dataBase.get("dataBaseKey").get("workKey"));
    }
}