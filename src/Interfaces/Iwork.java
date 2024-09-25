package Interfaces;

public interface Iwork {
    void add(Object key, Object value, int time);
    void delete(Object key);
    void put(Object key, Object value);
    void viewAll();
    void viewKey(Object key);
}
