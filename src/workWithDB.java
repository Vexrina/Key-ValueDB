import Interfaces.Iwork;
import java.util.concurrent.ConcurrentHashMap;

public class workWithDB implements Iwork {
    int time;
    ConcurrentHashMap<Object, Object> concurrentHashMap = new ConcurrentHashMap<>();
    /*
     * TO DO:
     * Добавить возможность пункта со временем, чтобы в дальнейшем удалять ключи (ttl)
     * */
    @Override
    public void add(Object key, Object value, int time) {
        this.concurrentHashMap.put(key, value);
    }

    @Override
    public void delete(Object key) {
        this.concurrentHashMap.remove(key);
    }

    /*
    * TO DO:
    * Если такого ключа нет в системе - добавить
    * Изменять время для ключа ttl
    * */
    @Override
    public void put(Object key, Object value) {
        this.concurrentHashMap.put(key, value);
    }

    @Override
    public void viewAll(){
        System.out.println(concurrentHashMap);
    }

    @Override
    public void viewKey(Object key) {
        System.out.println(concurrentHashMap.get(key));
    }
}
