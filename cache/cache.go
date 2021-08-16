package cache

import (
	"goffeine/cache/internal/fsketch"
	"goffeine/cache/internal/queue"
	"goffeine/cache/internal/node"
)

type Cache struct {
	cap        int
	sketch     *fsketch.FSketch
	windowQ    *queue.AccessOrderQueue
	probationQ *queue.AccessOrderQueue
	protectedQ *queue.AccessOrderQueue
}

func New(cap int) Cache {
	return Cache{
		cap:        cap,
		sketch:     fsketch.New(cap),
		windowQ:    queue.New(cap),
		probationQ: queue.New(cap),
		protectedQ: queue.New(cap),
	}
}

func (c *Cache) Capacity() int {
	return c.cap
}

func (c *Cache) Len() int {
	return c.windowQ.Len() + c.probationQ.Len() + c.protectedQ.Len()
}

func (c *Cache) Contains(key string) bool {
	return c.windowQ.Contains(key) || c.probationQ.Contains(key) || c.protectedQ.Contains(key)
}

// 往cache里面添加内容
func (c *Cache) Add(key string, value interface{}) {
	//pNode := node.New(key, value)

	// 如果不在cache里面，先添加到admission
	if !c.Contains(key) {
		//pNodeEliminated := c.windowQ.Push(key, pNode)
		//if pNodeEliminated != nil {
		//	return
		//}
		// 到这里表示admission满了，且自动淘汰了一个

	}
}

// Evicts entries if the cache exceeds the maximum.
//
//void evictEntries() {
//    if (!evicts()) {
//      return;
//    }
//    int candidates = evictFromWindow();
//    evictFromMain(candidates);
//  }
//
//

 /**
  * Evicts entries from the window space into the main space while the window size exceeds a
  * maximum.
  *
  * @return the number of candidate entries evicted from the window space
  */
 @GuardedBy("evictionLock")
 int evictFromWindow() {
   int candidates = 0;
   Node<K, V> node = accessOrderWindowDeque().peek();
   while (windowWeightedSize() > windowMaximum()) {
     // The pending operations will adjust the size to reflect the correct weight
     if (node == null) {
       break;
     }

     Node<K, V> next = node.getNextInAccessOrder();
     if (node.getPolicyWeight() != 0) {
       node.makeMainProbation();
       accessOrderWindowDeque().remove(node);
       accessOrderProbationDeque().add(node);
       candidates++;

       setWindowWeightedSize(windowWeightedSize() - node.getPolicyWeight());
     }
     node = next;
   }

   return candidates;
 }

//  /**
//   * Evicts entries from the main space if the cache exceeds the maximum capacity. The main space
//   * determines whether admitting an entry (coming from the window space) is preferable to retaining
//   * the eviction policy's victim. This is decision is made using a frequency filter so that the
//   * least frequently used entry is removed.
//   *
//   * The window space candidates were previously placed in the MRU position and the eviction
//   * policy's victim is at the AccessOrderQueue position. The two ends of the queue are evaluated while an
//   * eviction is required. The number of remaining candidates is provided and decremented on
//   * eviction, so that when there are no more candidates the victim is evicted.
//   *
//   * @param candidates the number of candidate entries evicted from the window space
//   */
//  @GuardedBy("evictionLock")
//  void evictFromMain(int candidates) {
//    int victimQueue = PROBATION;
//    Node<K, V> victim = accessOrderProbationDeque().peekFirst();
//    Node<K, V> candidate = accessOrderProbationDeque().peekLast();
//    while (weightedSize() > maximum()) {
//      // Stop trying to evict candidates and always prefer the victim
//      if (candidates == 0) {
//        candidate = null;
//      }
//
//      // Try evicting from the protectedQ and window queues
//      if ((candidate == null) && (victim == null)) {
//        if (victimQueue == PROBATION) {
//          victim = accessOrderProtectedDeque().peekFirst();
//          victimQueue = PROTECTED;
//          continue;
//        } else if (victimQueue == PROTECTED) {
//          victim = accessOrderWindowDeque().peekFirst();
//          victimQueue = WINDOW;
//          continue;
//        }
//
//        // The pending operations will adjust the size to reflect the correct weight
//        break;
//      }
//
//      // Skip over entries with zero weight
//      if ((victim != null) && (victim.getPolicyWeight() == 0)) {
//        victim = victim.getNextInAccessOrder();
//        continue;
//      } else if ((candidate != null) && (candidate.getPolicyWeight() == 0)) {
//        candidate = candidate.getPreviousInAccessOrder();
//        candidates--;
//        continue;
//      }
//
//      // Evict immediately if only one of the entries is present
//      if (victim == null) {
//        @SuppressWarnings("NullAway")
//        Node<K, V> previous = candidate.getPreviousInAccessOrder();
//        Node<K, V> evict = candidate;
//        candidate = previous;
//        candidates--;
//        evictEntry(evict, RemovalCause.SIZE, 0L);
//        continue;
//      } else if (candidate == null) {
//        Node<K, V> evict = victim;
//        victim = victim.getNextInAccessOrder();
//        evictEntry(evict, RemovalCause.SIZE, 0L);
//        continue;
//      }
//
//      // Evict immediately if an entry was collected
//      K victimKey = victim.getKey();
//      K candidateKey = candidate.getKey();
//      if (victimKey == null) {
//        @NonNull Node<K, V> evict = victim;
//        victim = victim.getNextInAccessOrder();
//        evictEntry(evict, RemovalCause.COLLECTED, 0L);
//        continue;
//      } else if (candidateKey == null) {
//        candidates--;
//        @NonNull Node<K, V> evict = candidate;
//        candidate = candidate.getPreviousInAccessOrder();
//        evictEntry(evict, RemovalCause.COLLECTED, 0L);
//        continue;
//      }
//
//      // Evict immediately if the candidate's weight exceeds the maximum
//      if (candidate.getPolicyWeight() > maximum()) {
//        candidates--;
//        Node<K, V> evict = candidate;
//        candidate = candidate.getPreviousInAccessOrder();
//        evictEntry(evict, RemovalCause.SIZE, 0L);
//        continue;
//      }
//
//      // Evict the entry with the lowest frequency
//      candidates--;
//      if (admit(candidateKey, victimKey)) {
//        Node<K, V> evict = victim;
//        victim = victim.getNextInAccessOrder();
//        evictEntry(evict, RemovalCause.SIZE, 0L);
//        candidate = candidate.getPreviousInAccessOrder();
//      } else {
//        Node<K, V> evict = candidate;
//        candidate = candidate.getPreviousInAccessOrder();
//        evictEntry(evict, RemovalCause.SIZE, 0L);
//      }
//    }
//  }
//
//  /**
//   * Determines if the candidate should be accepted into the main space, as determined by its
//   * frequency relative to the victim. A small amount of randomness is used to protect against hash
//   * collision attacks, where the victim's frequency is artificially raised so that no new entries
//   * are admitted.
//   *
//   * @param candidateKey the key for the entry being proposed for long term retention
//   * @param victimKey the key for the entry chosen by the eviction policy for replacement
//   * @return if the candidate should be admitted and the victim ejected
//   */
//  @GuardedBy("evictionLock")
//  boolean admit(K candidateKey, K victimKey) {
//    int victimFreq = frequencySketch().frequency(victimKey);
//    int candidateFreq = frequencySketch().frequency(candidateKey);
//    if (candidateFreq > victimFreq) {
//      return true;
//    } else if (candidateFreq <= 5) {
//      // The maximum frequency is 15 and halved to 7 after a reset to age the history. An attack
//      // exploits that a hot candidate is rejected in favor of a hot victim. The threshold of a warm
//      // candidate reduces the number of random acceptances to minimize the impact on the hit rate.
//      return false;
//    }
//    int random = ThreadLocalRandom.current().nextInt();
//    return ((random & 127) == 0);
//  }
//
