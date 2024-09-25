import main.workWithDB;

public class Main {
    public static void main(String[] args) {
        workWithDB db = new workWithDB();

        db.add(1, "Hello, world!", 1);

        System.out.println(db.viewAll());
    }
}