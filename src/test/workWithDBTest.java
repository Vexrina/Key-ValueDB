package test;
import static org.junit.jupiter.api.Assertions.*;
import main.workWithDB;

class workWithDBTest {
    workWithDB db = new workWithDB();
    @org.junit.jupiter.api.Test
    void add() {
        db.add(1, "Hello, world", 1);
        assertEquals(1, db.size());
    }

    @org.junit.jupiter.api.Test
    void delete() {
    }

    @org.junit.jupiter.api.Test
    void put() {
    }

    @org.junit.jupiter.api.Test
    void viewAll() {
    }

    @org.junit.jupiter.api.Test
    void viewKey() {
        db.add(1, "Hello, world", 1);
        assertEquals("Hello, world", db.viewValueByKey(1));
    }
}